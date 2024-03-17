package observer

import (
	"context"
	"fmt"
)

type SyncEventBus struct {
	*BaseEventBus
}

func NewSyncEventBus() *SyncEventBus {
	return &SyncEventBus{
		BaseEventBus: NewBaseEventBus(),
	}
}

func (s *SyncEventBus) Publish(ctx context.Context, e *Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	obs, ok := s.observers[e.Topic]
	if !ok {
		return
	}
	errs := make(map[Observer]error)
	for o := range obs {
		if err := o.OnChange(ctx, e); err != nil {
			errs[o] = err
		}
	}
	s.handleErr(ctx, errs)
}

func (s *SyncEventBus) handleErr(ctx context.Context, errs map[Observer]error) {
	for o, err := range errs {
		// ...处理逻辑
		fmt.Printf("observer: %v, err: %v\n", o, err)
	}
}
