package main

import (
	"context"
	"github.com/red-life/bus"
	"log"
)

type UserCreatedEvent struct {
	Username string
	FullName string
}

type UserCreatedEventHandler struct{}

func (UserCreatedEventHandler) Handle(ctx context.Context, event bus.Event) {
	userCreatedEvent := event.(UserCreatedEvent)
	log.Printf("New User Created! %+v\n", userCreatedEvent)
}

func main() {
	b := bus.NewBus[bus.Event]()
	eventBus := bus.NewEventBus(b)
	userCreatedEventHandler := new(UserCreatedEventHandler)
	err := eventBus.Subscribe("UserCreated", userCreatedEventHandler)
	if err != nil {
		panic(err)
	}
	eventBus.Publish(context.Background(), "UserCreated", UserCreatedEvent{
		Username: "Shayan",
		FullName: "Red Life",
	})
}
