package bus

import (
	"context"
	"errors"
)

type CommandHandler Handler
type Command Data

type ICommandBus interface {
	Subscribe(topic string, handler CommandHandler) error
	Unsubscribe(topic string, handler CommandHandler) error
	Subscriber(topic string) (CommandHandler, error)
	Publish(ctx context.Context, topic string, command Command)
	PublishAsync(ctx context.Context, topic string, command Command)
}

var (
	ErrTopicAlreadySubscribed = errors.New("this topic already been subscribed")
)

func NewCommandBus(bus *Bus) ICommandBus {
	return &CommandBus{
		bus: bus,
	}
}

type CommandBus struct {
	bus *Bus
}

func (c CommandBus) Subscribe(topic string, handler CommandHandler) error {
	_, err := c.bus.GetHandlers(topic)
	if errors.Is(err, ErrTopicNotFound) {
		return c.bus.RegisterHandler(topic, handler)
	}
	return ErrTopicAlreadySubscribed
}

func (c CommandBus) Unsubscribe(topic string, handler CommandHandler) error {
	return c.bus.RegisterHandler(topic, handler)
}

func (c CommandBus) Subscriber(topic string) (CommandHandler, error) {
	handlers, err := c.bus.GetHandlers(topic)
	if err != nil {
		return nil, err
	}
	return handlers[0], nil
}

func (c CommandBus) Publish(ctx context.Context, topic string, command Command) {
	c.bus.Broadcast(ctx, topic, command)
}

func (c CommandBus) PublishAsync(ctx context.Context, topic string, command Command) {
	c.bus.BroadcastAsync(ctx, topic, command)
}
