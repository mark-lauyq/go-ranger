package pubsub

import (
	"context"
	"time"
)

// Publisher ...
type Publisher interface {
	// Publish will publish a message with context.
	Publish(context.Context, string, string) error
}

// Subscriber ...
type Subscriber interface {
	// Start will return a channel of raw messages.
	Start() <-chan Message
	// Err will contain any errors returned from the consumer connection.
	Err() error
	// Stop will initiate a graceful shutdown of the subscriber connection.
	Stop() error
}

// Message ...
type Message interface {
	String() string
	ExtendDoneDeadline(time.Duration) error
	Done() error
	GetReceiveCount() (int, error)
	GetMessageId() string
}
