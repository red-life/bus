package bus

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

var (
	ErrHandlerNotFound  = errors.New("handler not found")
	ErrTopicNotFound    = errors.New("topic not found")
	ErrDuplicateHandler = errors.New("duplicate handler")
)

type Handler[T any] interface {
	Handle(context.Context, T)
}

type broadcastData[T any] struct {
	ctx   context.Context
	topic string
	data  T
	async bool
}

func NewBus[T any]() *Bus[T] {
	bus := &Bus[T]{
		lock:     sync.RWMutex{},
		handlers: map[string][]Handler[T]{},
		ch:       make(chan broadcastData[T]),
		wg:       sync.WaitGroup{},
	}
	return bus
}

type Bus[T any] struct {
	lock     sync.RWMutex
	handlers map[string][]Handler[T]
	ch       chan broadcastData[T]
	wg       sync.WaitGroup
}

func (b *Bus[T]) RegisterHandler(topic string, handler Handler[T]) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, err := b.findHandlerIdx(topic, handler); !errors.Is(err, ErrHandlerNotFound) && !errors.Is(err, ErrTopicNotFound) {
		return ErrDuplicateHandler
	}
	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}

func (b *Bus[T]) UnregisterHandler(topic string, handler Handler[T]) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	idx, err := b.findHandlerIdx(topic, handler)
	if err != nil {
		return err
	}
	b.handlers[topic] = append(b.handlers[topic][:idx], b.handlers[topic][idx+1:]...)
	return nil
}

func (b *Bus[T]) findHandlerIdx(topic string, handler Handler[T]) (int, error) {
	handlers, ok := b.handlers[topic]
	if !ok {
		return 0, ErrTopicNotFound
	}
	value := reflect.ValueOf(handler)
	for i, h := range handlers {
		if reflect.ValueOf(h).Type() == value.Type() &&
			reflect.ValueOf(h).Pointer() == value.Pointer() {
			return i, nil
		}
	}
	return 0, ErrHandlerNotFound
}

func (b *Bus[T]) Broadcast(ctx context.Context, topic string, data T) {
	b.doBroadcast(broadcastData[T]{
		ctx:   ctx,
		topic: topic,
		data:  data,
		async: false,
	})
}

func (b *Bus[T]) BroadcastAsync(ctx context.Context, topic string, data T) {
	b.doBroadcast(broadcastData[T]{
		ctx:   ctx,
		topic: topic,
		data:  data,
		async: true,
	})
}

func (b *Bus[T]) GetHandlers(topic string) ([]Handler[T], error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	handlers, ok := b.handlers[topic]
	if !ok {
		return nil, ErrTopicNotFound
	}
	return handlers, nil
}

func (b *Bus[T]) WaitAsync() {
	b.wg.Wait()
}

func (b *Bus[T]) doBroadcast(data broadcastData[T]) {
	handlers, err := b.GetHandlers(data.topic)
	if err != nil {
		return
	}
	for _, h := range handlers {
		if data.async {
			b.wg.Add(1)
			go func() {
				h.Handle(data.ctx, data.data)
				b.wg.Done()
			}()
		} else {
			h.Handle(data.ctx, data.data)
		}
	}
}
