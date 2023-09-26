package main

import (
	"context"
	"github.com/red-life/bus"
	"log"
)

type CreateUserCommand struct {
	Username string
	Password string
}

type CreateUserCommandHandler struct{}

func (CreateUserCommandHandler) Handle(ctx context.Context, cmd bus.Command) {
	createUserCmd := cmd.(CreateUserCommand)
	log.Printf("Command Received! %+v\n", createUserCmd)
}

func main() {
	b := bus.NewBus[bus.Command]()
	commandBus := bus.NewCommandBus(b)
	createUserCommandHandler := new(CreateUserCommandHandler)
	err := commandBus.Subscribe("CreateUser", createUserCommandHandler)
	if err != nil {
		panic(err)
	}
	commandBus.Publish(context.Background(), "CreateUser", CreateUserCommand{
		Username: "Shayan",
		Password: "123456",
	})
}
