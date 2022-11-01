package fastrand

import (
	_ "unsafe" // for linkname
)

//go:linkname Fastrand runtime.fastrand
func Fastrand() uint32

// Uint32 returns a pseudo-random 32-bit value as a uint32.
var Uint32 = Fastrand

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func Uint64() uint64 {
	return (uint64(Fastrand()) << 32) | uint64(Fastrand())
}
func Int63() int64 {
	// EQ
	return int64(Uint64() & (1<<63 - 1))
}

func Int63n(n int64) int64 {
	// EQ
	if n <= 0 {
		panic("invalid argument to Int63n")
	}
	if n&(n-1) == 0 { // n is power of two, can mask
		return Int63() & (n - 1)
	}
	max := int64((1 << 63) - 1 - (1<<63)%uint64(n))
	v := Int63()
	for v > max {
		v = Int63()
	}
	return v % n
}
func Intn(n int) int {
	// EQ
	if n <= 0 {
		panic("invalid argument to Intn")
	}
	if n <= 1<<31-1 {
		return int(Int31n(int32(n)))
	}
	return int(Int63n(int64(n)))
}
func Int31n(n int32) int32 {
	// EQ
	if n <= 0 {
		panic("invalid argument to Int31n")
	}
	v := Uint32()
	prod := uint64(v) * uint64(n)
	low := uint32(prod)
	if low < uint32(n) {
		thresh := uint32(-n) % uint32(n)
		for low < thresh {
			v = Uint32()
			prod = uint64(v) * uint64(n)
			low = uint32(prod)
		}
	}
	return int32(prod >> 32)
}
