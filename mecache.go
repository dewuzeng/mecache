package mecache

import (
	"runtime"
	"sync"
	"time"
)

func assertCacheImplementation() {
	var _ Cache = (*MeCache)(nil)
}

const DefaultExpiration time.Duration = 0

type MeCache struct {
	defaultExpiration time.Duration
	items             map[string]Item
	mu                sync.RWMutex
	monitor           *monitor
}

func (mc *MeCache) Set(k string, x interface{}, d time.Duration) {
	if d == DefaultExpiration {
		d = mc.defaultExpiration
	}
	var exp int64
	if d > 0 {
		exp = time.Now().Add(d).UnixNano()
	}
	mc.mu.Lock()
	mc.items[k] = Item{
		Val: x,
		Exp: exp,
	}
	mc.mu.Unlock()
}

func (mc *MeCache) SetDefault(k string, x interface{}) {
	mc.Set(k, x, DefaultExpiration)
}

func (mc *MeCache) Get(k string) (interface{}, bool) {
	mc.mu.RLock()
	item, found := mc.items[k]
	mc.mu.RUnlock()
	if !found || (item.Exp > 0 && time.Now().UnixNano() > item.Exp) {
		return nil, false
	}
	return item.Val, true
}

func (mc *MeCache) DeleteExpired() {
	now := time.Now().UnixNano()
	mc.mu.Lock()
	for k, v := range mc.items {
		if v.Exp > 0 && now > v.Exp {
			delete(mc.items, k)
		}
	}
	mc.mu.Unlock()
}

func New(defaultExpiration, cleanupInterval time.Duration) *MeCache {
	if defaultExpiration <= 0 {
		defaultExpiration = DefaultExpiration
	}
	items := make(map[string]Item)
	mc := &MeCache{
		defaultExpiration: defaultExpiration,
		items:             items,
	}
	if cleanupInterval > 0 {
		startMonitor(mc, cleanupInterval)
		runtime.SetFinalizer(mc, stopMonitor)
	}
	return mc
}

type Item struct {
	Val interface{}
	Exp int64
}

type monitor struct {
	Interval time.Duration
	stop     chan bool
}

func (m *monitor) Run(mc *MeCache) {
	ticker := time.NewTicker(m.Interval)
	for {
		select {
		case <-ticker.C:
			mc.DeleteExpired()
		case <-m.stop:
			ticker.Stop()
			return
		}
	}
}

func startMonitor(mc *MeCache, ci time.Duration) {
	m := &monitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	mc.monitor = m
	go m.Run(mc)
}

func stopMonitor(mc *MeCache) {
	mc.monitor.stop <- true
}
