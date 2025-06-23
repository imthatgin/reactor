package reactor

import (
	"sync"
)

type Subscriber interface {
	Notify()
}

type Signal[T comparable] struct {
	mu          sync.RWMutex
	value       T
	subscribers map[Subscriber]struct{}
}

func NewSignal[T comparable](value T) *Signal[T] {
	return &Signal[T]{
		value:       value,
		subscribers: make(map[Subscriber]struct{}),
	}
}

func (s *Signal[T]) Get() T {
	//s.mu.RLock()
	//defer s.mu.RUnlock()

	if tr := currentTracker(); tr != nil {
		tr.addDependency(s)
	}

	return s.value
}

func (s *Signal[T]) Set(newValue T) {
	//s.mu.Lock()
	//defer s.mu.Unlock()

	// Check for equality to avoid needlessly updating
	if s.value == newValue {
		return
	}

	s.value = newValue
	for sub := range s.subscribers {
		sub.Notify()
	}
}

func (s *Signal[T]) Subscribe(sub Subscriber) {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	s.subscribers[sub] = struct{}{}
}

func (s *Signal[T]) Unsubscribe(sub Subscriber) {
	//s.mu.Lock()
	//defer s.mu.Unlock()
	delete(s.subscribers, sub)
}
