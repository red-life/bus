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

type Data any

type Handler interface {
	Handle(context.Context, Data)
}

type broadcastData struct {
	ctx   context.Context
	topic string
	data  Data
	async bool
}

func NewBus() *Bus {
	bus := &Bus{
		lock:     sync.RWMutex{},
		handlers: map[string][]Handler{},
		ch:       make(chan broadcastData),
		wg:       sync.WaitGroup{},
	}
	return bus
}

type Bus struct {
	lock     sync.RWMutex
	handlers map[string][]Handler
	ch       chan broadcastData
	wg       sync.WaitGroup
}

func (b *Bus) RegisterHandler(topic string, handler Handler) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	if _, err := b.findHandlerIdx(topic, handler); !errors.Is(err, ErrHandlerNotFound) && !errors.Is(err, ErrTopicNotFound) {
		return ErrDuplicateHandler
	}
	b.handlers[topic] = append(b.handlers[topic], handler)
	return nil
}

func (b *Bus) UnregisterHandler(topic string, handler Handler) error {
	b.lock.Lock()
	defer b.lock.Unlock()
	idx, err := b.findHandlerIdx(topic, handler)
	if err != nil {
		return err
	}
	b.handlers[topic] = append(b.handlers[topic][:idx], b.handlers[topic][idx+1:]...)
	return nil
}

func (b *Bus) findHandlerIdx(topic string, handler Handler) (int, error) {
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

func (b *Bus) Broadcast(ctx context.Context, topic string, data Data) {
	b.doBroadcast(broadcastData{
		ctx:   ctx,
		topic: topic,
		data:  data,
		async: false,
	})
}

func (b *Bus) BroadcastAsync(ctx context.Context, topic string, data Data) {
	b.doBroadcast(broadcastData{
		ctx:   ctx,
		topic: topic,
		data:  data,
		async: true,
	})
}

func (b *Bus) GetHandlers(topic string) ([]Handler, error) {
	b.lock.RLock()
	defer b.lock.RUnlock()
	handlers, ok := b.handlers[topic]
	if !ok {
		return nil, ErrTopicNotFound
	}
	return handlers, nil
}

func (b *Bus) WaitAsync() {
	b.wg.Wait()
}

func (b *Bus) doBroadcast(data broadcastData) {
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
