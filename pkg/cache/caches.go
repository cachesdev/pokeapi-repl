package cache

import (
	"sync"
	"time"
)

type Cache struct {
	kv           *sync.Map
	reapInterval time.Duration
}

func NewCache(reapInterval time.Duration) *Cache {
	c := &Cache{
		kv:           &sync.Map{},
		reapInterval: reapInterval,
	}

	go func() {
		ticker := time.NewTicker(c.reapInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.reap()
			}
		}
	}()

	return c
}

func (c *Cache) reap() {
	c.kv.Range(func(key, value any) bool {
		obj := value.(object)
		if time.Since(obj.createdAt) > c.reapInterval {
			c.kv.Delete(key)
		}

		return true
	})
}

type object struct {
	createdAt time.Time
	val       any
}

func (c *Cache) Set(key any, val any) {
	obj := object{
		createdAt: time.Now(),
		val:       val,
	}

	c.kv.Store(key, obj)
}

func (c *Cache) Get(key any) (any, bool) {
	val, ok := c.kv.Load(key)
	if !ok {
		return nil, ok
	}

	obj := val.(object)

	return obj.val, ok
}
