package main

import (
	"sync"
)

type SafeMap struct {
	mu sync.Mutex
	m  map[string]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		m: make(map[string]interface{}),
	}
}

func (sm *SafeMap) Set(key string, value interface{}) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.m[key] = value
}

func (sm *SafeMap) Get(key string) (interface{}, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	val, ok := sm.m[key]
	return val, ok
}

func main() {
	// Реализация с помошью mutex, все действия потокобезопасны
	safeMap := NewSafeMap()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + string(rune(i))
			safeMap.Set(key, i)
		}(i)
	}

	wg.Wait()
}
