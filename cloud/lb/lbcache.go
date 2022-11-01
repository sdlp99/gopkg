package lb

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"sync"
	"sync/atomic"
	"time"
)

type Request struct {
}

type cacheResult struct {
	res         atomic.Value // newest and previous discovery result
	expire      int32        // 0 = normal, 1 = expire and collect next ticker
	serviceName string       // service psm
}

var (
	balancerFactories    sync.Map // key: resolver name + load-balancer name
	balancerFactoriesSfg singleflight.Group
)

func cacheKey(resolver, balancer string, opts Options) string {
	return fmt.Sprintf("%s|%s|{%s %s}", resolver, balancer, opts.RefreshInterval, opts.ExpireInterval)
}

type BalancerFactory struct {
	opts     Options
	cache    sync.Map // key -> LoadBalancer
	resolver Resolver
	balancer Loadbalancer
	sfg      singleflight.Group
}

type Config struct {
	Resolver Resolver
	Balancer Loadbalancer
	LbOpts   Options
}

// NewBalancerFactory get or create a balancer with given target.
// If it has the same key(resolver.Target(target)), we will cache and reuse the Balance.
func NewBalancerFactory(config Config) *BalancerFactory {
	config.LbOpts.Check()
	uniqueKey := cacheKey(config.Resolver.Name(), config.Balancer.Name(), config.LbOpts)
	val, ok := balancerFactories.Load(uniqueKey)
	if ok {
		return val.(*BalancerFactory)
	}
	val, _, _ = balancerFactoriesSfg.Do(uniqueKey, func() (interface{}, error) {
		b := &BalancerFactory{
			opts:     config.LbOpts,
			resolver: config.Resolver,
			balancer: config.Balancer,
		}
		go b.watcher()
		go b.refresh()
		balancerFactories.Store(uniqueKey, b)
		return b, nil
	})
	return val.(*BalancerFactory)
}

// watch expired balancer
func (b *BalancerFactory) watcher() {
	for range time.Tick(b.opts.ExpireInterval) {
		b.cache.Range(func(key, value interface{}) bool {
			cache := value.(*cacheResult)
			if atomic.CompareAndSwapInt32(&cache.expire, 0, 1) {
				// 1. set expire flag
				// 2. wait next ticker for collect, maybe the balancer is used again
				// (avoid being immediate delete the balancer which had been created recently)
			} else {
				b.cache.Delete(key)
				b.balancer.Delete(key.(string))
			}
			return true
		})
	}
}

// cache key with resolver name prefix avoid conflict for balancer
func renameResultCacheKey(res *Result, resolverName string) {
	res.CacheKey = resolverName + ":" + res.CacheKey
}

// refresh is used to update service discovery information periodically.
func (b *BalancerFactory) refresh() {
	for range time.Tick(b.opts.RefreshInterval) {
		b.cache.Range(func(key, value interface{}) bool {
			res, err := b.resolver.Resolve(context.Background(), key.(string))
			if err != nil {
				//hlog.Warnf("Hertz: resolver refresh failed, key=%s error=%s", key, err.Error())
				return true
			}
			renameResultCacheKey(&res, b.resolver.Name())
			cache := value.(*cacheResult)
			cache.res.Store(res)
			atomic.StoreInt32(&cache.expire, 0)
			b.balancer.Rebalance(res)
			return true
		})
	}
}

func (b *BalancerFactory) GetInstance(ctx context.Context, req string, tag map[string]string) (Instance, error) {
	cacheRes, err := b.getCacheResult(ctx, req, tag)
	if err != nil {
		return nil, err
	}
	atomic.StoreInt32(&cacheRes.expire, 0)
	ins := b.balancer.Pick(cacheRes.res.Load().(Result))
	if ins == nil {
		//hlog.Errorf("HERTZ: null instance. serviceName: %s, options: %v", string(req.Host()), req.Options())
		return nil, errors.New("instance not found")
	}
	return ins, nil
}

func (b *BalancerFactory) getCacheResult(ctx context.Context, host string, tag map[string]string) (*cacheResult, error) {
	target := b.resolver.Target(ctx, &TargetInfo{Host: string(host), Tags: tag})
	cr, existed := b.cache.Load(target)
	if existed {
		return cr.(*cacheResult), nil
	}
	cr, err, _ := b.sfg.Do(target, func() (interface{}, error) {
		cache := &cacheResult{
			serviceName: string(host),
		}
		res, err := b.resolver.Resolve(ctx, target)
		if err != nil {
			return cache, err
		}
		renameResultCacheKey(&res, b.resolver.Name())
		cache.res.Store(res)
		atomic.StoreInt32(&cache.expire, 0)
		b.balancer.Rebalance(res)
		b.cache.Store(target, cache)
		return cache, nil
	})
	if err != nil {
		return nil, err
	}
	return cr.(*cacheResult), nil
}
