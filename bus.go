package bus

import (
	"context"
	"errors"
	"reflect"
	"sync"
)

var (
	ErrHandlerNotFound = errors.New("handler not found")
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
		lock: sync.RWMutex{},
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
}

func (b *Bus) RegisterHandler(topic string, handler Handler) error {
	b.lock.Lock()
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
	b.lock.Unlock() // Not using defer due to performance issues
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

func (b *Bus) Broadcast(topic string, data Data) {
	b.ch <- broadcastData{
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
					go h.Handle(broadcastedData.ctx, broadcastedData.data)
				} else {
					h.Handle(broadcastedData.ctx, broadcastedData.data)
				}
			}
		}
	}
}
