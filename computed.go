package reactor

import (
	"sync"
)

type ComputeFunc[T any] func() T

type Computed[T any] struct {
	mu           sync.Mutex
	compute      ComputeFunc[T]
	value        T
	dirty        bool
	subscribers  map[Subscriber]struct{}
	dependencies map[subscribable]struct{}
}

func NewComputed[T any](compute ComputeFunc[T]) *Computed[T] {
	return &Computed[T]{
		compute:     compute,
		dirty:       true,
		subscribers: make(map[Subscriber]struct{}),
	}
}

func (c *Computed[T]) Get() T {
	//c.mu.Lock()
	//defer c.mu.Unlock()

	if c.dirty {
		c.recompute()
	}
	return c.value
}

func (c *Computed[T]) recompute() {
	// Unsubscribe from old deps
	for dep := range c.dependencies {
		dep.Unsubscribe(c)
	}
	c.dependencies = make(map[subscribable]struct{})

	pushTracker(c)
	c.value = c.compute()
	popTracker()
	c.dirty = false
}

func (c *Computed[T]) Notify() {
	//c.mu.Lock()
	//defer c.mu.Unlock()
	c.dirty = true
	for sub := range c.subscribers {
		sub.Notify()
	}
}

func (c *Computed[T]) addDependency(dep subscribable) {
	if _, ok := c.dependencies[dep]; !ok {
		c.dependencies[dep] = struct{}{}
		dep.Subscribe(c)
	}
}

func (c *Computed[T]) Subscribe(sub Subscriber) {
	//c.mu.Lock()
	//defer c.mu.Unlock()
	c.subscribers[sub] = struct{}{}
}

func (c *Computed[T]) Unsubscribe(sub Subscriber) {
	//c.mu.Lock()
	//defer c.mu.Unlock()
	delete(c.subscribers, sub)
}
