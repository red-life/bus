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

func NewBus(workers uint) *Bus {
	bus := &Bus{
		lock:     sync.RWMutex{},
		handlers: map[string][]Handler{},
		ch:       make(chan broadcastData, 0),
		wg:       sync.WaitGroup{},
	}
	for i := 0; i < int(workers); i++ {
		go bus.broadcaster()
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
	if idx := b.findHandlerIdx(topic, handler); idx != -1 {
		return ErrDuplicateHandler
	}
	b.handlers[topic] = append(b.handlers[topic], handler)
	b.lock.Unlock() // Not using defer due to performance issues
	return nil
}

func (b *Bus) UnregisterHandler(topic string, handler Handler) error {
	b.lock.Lock()
	idx := b.findHandlerIdx(topic, handler)
	if idx == -1 {
		return ErrHandlerNotFound
	}
	b.handlers[topic] = append(b.handlers[topic][:idx], b.handlers[topic][idx+1:]...)
	b.lock.Unlock()
	return nil
}

func (b *Bus) findHandlerIdx(topic string, handler Handler) int {
	handlers, ok := b.handlers[topic]
	if !ok {
		return -1
	}
	value := reflect.ValueOf(handler)
	for i, h := range handlers {
		if reflect.ValueOf(h).Type() == value.Type() &&
			reflect.ValueOf(h).Pointer() == value.Pointer() {
			return i
		}
	}
	return -1
}

func (b *Bus) Broadcast(ctx context.Context, topic string, data Data) {
	b.ch <- broadcastData{
		ctx:   ctx,
		topic: topic,
		data:  data,
		async: false,
	}
}

func (b *Bus) BroadcastAsync(ctx context.Context, topic string, data Data) {
	b.ch <- broadcastData{
		ctx:   ctx,
		topic: topic,
		data:  data,
		async: true,
	}
}

func (b *Bus) GetHandlers(topic string) ([]Handler, error) {
	b.lock.RLock()
	handlers, ok := b.handlers[topic]
	if !ok {
		return nil, ErrTopicNotFound
	}
	b.lock.RUnlock()
	return handlers, nil
}

func (b *Bus) WaitAsync() {
	b.wg.Wait()
}

func (b *Bus) broadcaster() {
	for {
		select {
		case broadcastedData := <-b.ch:
			b.lock.RLock()
			handlers, ok := b.handlers[broadcastedData.topic]
			b.lock.RUnlock()
			if !ok {
				continue
			}
			for _, h := range handlers {
				if broadcastedData.async {
					b.wg.Add(1)
					go func() {
						h.Handle(broadcastedData.ctx, broadcastedData.data)
						b.wg.Done()
					}()
				} else {
					h.Handle(broadcastedData.ctx, broadcastedData.data)
				}
			}
		}
	}
}
