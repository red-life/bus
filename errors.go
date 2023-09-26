package bus

import "errors"

// Bus errors (CommandBus & EventBus can also return these errors)
var (
	ErrHandlerNotFound  = errors.New("handler not found")
	ErrTopicNotFound    = errors.New("topic not found")
	ErrDuplicateHandler = errors.New("duplicate handler")
)

// CommandBus errors
var (
	ErrTopicAlreadySubscribed = errors.New("this topic already been subscribed")
)
