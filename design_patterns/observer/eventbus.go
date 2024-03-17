package observer

import (
	"context"
	"sync"
)

type BaseEventBus struct {
	mu        sync.RWMutex
	observers map[string]map[Observer]struct{}
}

func NewBaseEventBus() *BaseEventBus {
	return &BaseEventBus{
		mu:        sync.RWMutex{},
		observers: make(map[string]map[Observer]struct{}),
	}
}

func (b *BaseEventBus) Subscribe(topic string, o Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.observers[topic]
	if !ok {
		b.observers[topic] = make(map[Observer]struct{})
	}
	b.observers[topic][o] = struct{}{}
}

func (b *BaseEventBus) UnSubscribe(topic string, o Observer) {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.observers[topic]
	if !ok {
		return
	}
	_, ok = b.observers[topic][o]
	if !ok {
		return
	}
	delete(b.observers[topic], o)
}

func (b *BaseEventBus) Publish(ctx context.Context, e *Event) {

}
