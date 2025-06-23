package reactor

import (
	"sync"
)

type dependencyTracker interface {
	addDependency(dep subscribable)
}

type subscribable interface {
	Subscribe(Subscriber)
	Unsubscribe(Subscriber)
}

var (
	trackerMutex sync.Mutex
	trackerStack []dependencyTracker
)

func pushTracker(t dependencyTracker) {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()
	trackerStack = append(trackerStack, t)
}

func popTracker() {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()
	if len(trackerStack) > 0 {
		trackerStack = trackerStack[:len(trackerStack)-1]
	}
}

func currentTracker() dependencyTracker {
	trackerMutex.Lock()
	defer trackerMutex.Unlock()
	if len(trackerStack) == 0 {
		return nil
	}
	return trackerStack[len(trackerStack)-1]
}
