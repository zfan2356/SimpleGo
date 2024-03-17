package observer

import (
	"context"
	"testing"
)

func TestSyncEventBus(t *testing.T) {
	obA := NewBaseObserver("A")
	obB := NewBaseObserver("B")
	obC := NewBaseObserver("C")
	obD := NewBaseObserver("D")

	subs := NewSyncEventBus()
	topic := "testSync"
	subs.Subscribe(topic, obA)
	subs.Subscribe(topic, obB)
	subs.Subscribe(topic, obC)
	subs.Subscribe(topic, obD)

	subs.Publish(context.TODO(), &Event{
		Topic: topic,
		Val:   "test_publish",
	})
}
func TestAsyncEventBus(t *testing.T) {
	obA := NewBaseObserver("A")
	obB := NewBaseObserver("B")
	obC := NewBaseObserver("C")
	obD := NewBaseObserver("D")

	subs := NewAsyncEventBus()
	topic := "testSync"
	subs.Subscribe(topic, obA)
	subs.Subscribe(topic, obB)
	subs.Subscribe(topic, obC)
	subs.Subscribe(topic, obD)

	subs.Publish(context.TODO(), &Event{
		Topic: topic,
		Val:   "test_publish",
	})
}
