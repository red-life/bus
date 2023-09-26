# How to use EventBus
First of all, you have to create an instance of `Bus` because EventBus uses `Bus` to handle publishing and subscribing
```go
b := bus.NewBus[bus.Event]()
```
Then, easily pass the `b` to `NewEventBus` to create an instance of EventBus
```go
eventBus := bus.NewEventBus(b)
```
That's it. You have an Event Bus that you can easily publish and subscribe to different topics.\
*`NewEventBus` returns an interface. You can see the available methods at `IEventBus` interface*

### Subscribing to a topic
`Subscribe` method takes two arguments, the topic name and event handler. Handler is just an interface.
```go
type eventHandler struct {}
func (eventHandler) Handle(ctx context.Context, event bus.Event) {
	fmt.Printf("An event received!\n%+v\n", event)
}

handler := &eventHandler{}
eventBus.Subscribe("myTopic", handler)
```

### Publishing to a topic
There are two methods for publishing:
- `Publish`
- `PublishAsync`\
The key difference between them is `PublishAsync` creates a goroutine foreach subscriber but `Publish` waits foreach subscriber to finish its handling of the event.

```go
type myEvent struct {
	message string
	users []int
}
event := myEvent{
	message: "event published",
	users: []int{0, 1, 2},
}
eventBus.Publish("myTopic", event)
```

### Error Handling
All errors that might be returned are available in [errors.go](../errors.go)


### Example
You can the full example in [event_bus.go](../examples/event_bus.go)\
Also you can run the example by running the command below:\
`go run examples/event_bus.go`

