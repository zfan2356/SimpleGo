package observer

import "context"

type Event struct {
	Topic string
	Val   interface{}
}

type Observer interface {
	OnChange(ctx context.Context, e *Event) error
}

type EventBus interface {
	Subscribe(topic string, o Observer)
	UnSubscribe(topic string, o Observer)
	Publish(ctx context.Context, e *Event)
}
