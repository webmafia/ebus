package ebus

import "context"

// Convenient method for subscribing for an event and pushling to a channel, so that
// slow subscribers can do their work in a goroutine. If channel is full, any publisher
// will block.
// Publishers must publish the specific variable type.
func SubToChan[T any](eb *EventBus, event Event, ch chan<- T) func(context.Context, *T) {
	return Sub(eb, event, func(_ context.Context, val *T) {
		// TODO: Respect context.Context?
		ch <- *val
	})
}
