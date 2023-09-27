# How to use CommandBus
First of all, you have to create an instance of `Bus` because CommandBUs uses `Bus` to handle publishing and subscribing
```go
b := bus.NewBus[bus.Command]()
```
Then, easily pass the `b` to `NewCommandBus` to create an instance of CommandBus
```go
commandBus := bus.NewCommandBus(b)
```
That's it. You have a Command Bus that you can easily use in your CQRS Architecture.\
*`NewCommandBus` returns an interface. You can see the available methods at `ICommandBus` interface*

### Subscribing to a command
The `Subscribe` method takes two arguments, the topic name and the command handler. The handler is just an interface.
```go
type commandHandler struct {}
func (commandHandler) Handle(ctx context.Context, cmd bus.Command) {
	fmt.Printf("A command received!\n%+v\n", cmd)
}

handler := &commandHandler{}
commandBus.Subscribe("myTopic", handler)
```

### Publishing to a topic
There are two methods for publishing:
- `Publish`
- `PublishAsync`\
  The key difference between them is that `PublishAsync` creates a goroutine for each subscriber but `Publish` waits for the subscriber to finish its handling of the command.

```go
type myCommand struct {
	username string
	password string
}
cmd := myCommand{
	username: "shayan",
    password: "123456",
}
commandBus.Publish("myCommand", cmd)
```

### Error Handling
All errors that might be returned are available in [errors.go](../errors.go)


### Example
You can the full example in [command_bus.go](../examples/command_bus.go)\
Also, you can run the example by running the command below:\
`go run examples/command_bus.go`

