package observer

import (
	"context"
	"fmt"
)

/*
异步模式下的EventBus, 可以采用异步启动一个守护协程,然后对接收到
的错误,进行异步处理
*/
type AsyncEventBus struct {
	*BaseEventBus
	errChan chan *observerWithErr
	ctx     context.Context
	stop    context.CancelFunc
}

type observerWithErr struct {
	o   Observer
	err error
}

func NewAsyncEventBus() *AsyncEventBus {
	aBus := AsyncEventBus{
		BaseEventBus: NewBaseEventBus(),
		errChan:      make(chan *observerWithErr),
	}
	aBus.ctx, aBus.stop = context.WithCancel(context.Background())
	go aBus.handleErr()
	return &aBus
}

func (a *AsyncEventBus) Stop() {
	a.stop()
}

func (a *AsyncEventBus) handleErr() {
	for {
		select {
		case <-a.ctx.Done():
			return
		case resp := <-a.errChan:
			// ..处理逻辑
			fmt.Printf("observer: % v, err: %v", resp.o, resp.err)
		}
	}
}

func (a *AsyncEventBus) Publish(ctx context.Context, e *Event) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	obs, ok := a.observers[e.Topic]
	if !ok {
		return
	}

	for o := range obs {
		o := o
		go func() {
			if err := o.OnChange(ctx, e); err != nil {
				select {
				case <-a.ctx.Done():
				case a.errChan <- &observerWithErr{o, err}:
				}
			}
		}()
	}
}
