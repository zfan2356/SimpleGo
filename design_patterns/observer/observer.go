package observer

import (
	"context"
	"fmt"
)

type BaseObserver struct {
	name string
}

func NewBaseObserver(name string) *BaseObserver {
	return &BaseObserver{
		name: name,
	}
}

func (b *BaseObserver) OnChange(ctx context.Context, e *Event) error {
	fmt.Printf("observer: %s, event key: %s, event val: %v\n", b.name, e.Topic, e.Val)
	// ...
	return nil
}
