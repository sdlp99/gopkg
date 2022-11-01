package lb

import (
	"github.com/sdlp99/sdpkg/internal/fastrand"
	"golang.org/x/sync/singleflight"
	"sync"
)

type weightedBalancer struct {
	cachedWeightInfo sync.Map
	sfg              singleflight.Group
}

type weightInfo struct {
	instances []Instance
	entries   []int
	weightSum int
}

func NewWeightedBalancer() Loadbalancer {
	lb := &weightedBalancer{}
	return lb
}

func (wb *weightedBalancer) calcWeightInfo(e Result) *weightInfo {
	w := &weightInfo{
		instances: make([]Instance, len(e.Instances)),
		weightSum: 0,
		entries:   make([]int, len(e.Instances)),
	}

	var cnt int

	for idx := range e.Instances {
		weight := e.Instances[idx].Weight()
		if weight > 0 {
			w.instances[cnt] = e.Instances[idx]
			w.entries[cnt] = weight
			w.weightSum += weight
			cnt++
		} else {
			//hlog.Warnf("HERTZ: Invalid weight=%d on instance address=%s", weight, e.Instances[idx].Address())
		}
	}

	w.instances = w.instances[:cnt]

	return w
}

// Pick implements the Loadbalancer interface.
func (wb *weightedBalancer) Pick(e Result) Instance {
	wi, ok := wb.cachedWeightInfo.Load(e.CacheKey)
	if !ok {
		wi, _, _ = wb.sfg.Do(e.CacheKey, func() (interface{}, error) {
			return wb.calcWeightInfo(e), nil
		})
		wb.cachedWeightInfo.Store(e.CacheKey, wi)
	}

	w := wi.(*weightInfo)
	if w.weightSum <= 0 {
		return nil
	}

	weight := fastrand.Intn(w.weightSum)
	for i := 0; i < len(w.instances); i++ {
		weight -= w.entries[i]
		if weight < 0 {
			return w.instances[i]
		}
	}

	return nil
}

// Rebalance implements the Loadbalancer interface.
func (wb *weightedBalancer) Rebalance(e Result) {
	wb.cachedWeightInfo.Store(e.CacheKey, wb.calcWeightInfo(e))
}

// Delete implements the Loadbalancer interface.
func (wb *weightedBalancer) Delete(cacheKey string) {
	wb.cachedWeightInfo.Delete(cacheKey)
}

func (wb *weightedBalancer) Name() string {
	return "weight_random"
}
