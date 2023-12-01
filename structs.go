package main

import (
	"sync"
	"sync/atomic"
)

type SliceWithCount struct {
	sync.Mutex
	value [][]int
}

func (slice *SliceWithCount) GetSlice() [][]int {
	slice.Lock()
	defer slice.Unlock()

	return slice.value
}

func (slice *SliceWithCount) AddToSlice(x, y int) {
	slice.Lock()
	defer slice.Unlock()

	if slice.value[y][x] != 0 {
		return
	} else {
		atomic.AddInt64(&UniqueMaxCount, 1)
		slice.value[y][x]++
	}
}

type Dot struct {
	x int
	y int
}

type SetWithDots struct {
	sync.Mutex
	value map[Dot]struct{}
}

func (set *SetWithDots) AddDot(x, y int) {
	set.Lock()
	defer set.Unlock()

	if _, ok := set.value[Dot{x, y}]; ok {
		return
	} else {
		atomic.AddInt64(&UniqueMaxCount, 1)
		set.value[Dot{x, y}] = struct{}{}
	}
}

func (set *SetWithDots) GetSet() map[Dot]struct{} {
	set.Lock()
	defer set.Unlock()

	return set.value
}

func (set *SetWithDots) Clear() {
	set.Lock()
	defer set.Unlock()

	UniqueMaxCount = 0

	for key := range set.value {
		delete(set.value, key)
	}
}
