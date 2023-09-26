package bus

import "context"

type Event Data
type EventHandler Handler

type IEventBus interface {
	Subscribe(topic string, handler EventHandler) error
	Unsubscribe(topic string, handler EventHandler) error
	Subscribers(topic string) ([]EventHandler, error)
	Publish(ctx context.Context, topic string, event Event)
	PublishAsync(ctx context.Context, topic string, event Event)
}

func NewEventBus(bus *Bus) IEventBus {
	return &EventBus{
		bus: bus,
	}
}

type EventBus struct {
	bus *Bus
}

func (e EventBus) Subscribe(topic string, handler EventHandler) error {
	return e.bus.RegisterHandler(topic, handler)
}

func (e EventBus) Unsubscribe(topic string, handler EventHandler) error {
	return e.bus.RegisterHandler(topic, handler)
}

func (e EventBus) Subscribers(topic string) ([]EventHandler, error) {
	handlers, err := e.bus.GetHandlers(topic)
	if err != nil {
		return nil, err
	}
	eventHandlers := make([]EventHandler, 0, len(handlers))
	for _, h := range handlers {
		eventHandlers = append(eventHandlers, EventHandler(h))
	}
	return eventHandlers, nil
}

func (e EventBus) Publish(ctx context.Context, topic string, event Event) {
	e.bus.Broadcast(ctx, topic, event)
}

func (e EventBus) PublishAsync(ctx context.Context, topic string, event Event) {
	e.bus.BroadcastAsync(ctx, topic, event)
}
