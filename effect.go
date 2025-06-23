package reactor

import (
	"sync"
)

type EffectFunc func()

type Effect struct {
	mu      sync.Mutex
	fn      EffectFunc
	running bool
	deps    map[subscribable]struct{}
	self    Subscriber
}

func NewEffect(fn EffectFunc) *Effect {
	e := &Effect{
		fn:   fn,
		deps: make(map[subscribable]struct{}),
	}
	e.self = e
	e.run()
	return e
}

func (e *Effect) run() {
	//e.mu.Lock()
	if e.running {
		//e.mu.Unlock()
		return
	}
	e.running = true
	// Unsubscribe from old deps
	for dep := range e.deps {
		dep.Unsubscribe(e.self)
	}
	e.deps = make(map[subscribable]struct{})
	//e.mu.Unlock()

	pushTracker(e)
	e.fn()
	popTracker()

	//e.mu.Lock()
	e.running = false
	//e.mu.Unlock()
}

func (e *Effect) addDependency(dep subscribable) {
	//e.mu.Lock()
	//defer e.mu.Unlock()
	if _, ok := e.deps[dep]; !ok {
		e.deps[dep] = struct{}{}
		dep.Subscribe(e.self)
	}
}

func (e *Effect) Notify() {
	e.run()
}
