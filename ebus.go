package ebus

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/webmafia/ebus/list"
)

type Event uint32

// A thread-safe event bus.
// The zero EventBus is empty and ready for use.
// An EventBus must not be copied after first use.
//
//	bus := NewEventBus()
//
//	// Either send only event
//	bus.Pub(ctx, 123)
//
//	// Or also send variable
//	ebus.Pub(bus, ctx, 123, &myVar)
type EventBus struct {
	cbs sync.Map
}

type eventKey struct {
	event Event
	typ   uint32
}

// Convenient function to create a new event bus.
func NewEventBus() *EventBus {
	return &EventBus{}
}

// Returns the number of current subscribers
func (eb *EventBus) Subscribers() (n int64) {
	eb.cbs.Range(func(_, val any) bool {
		ptr := (*atomic.Int64)(efaceData(val))
		n += ptr.Load()
		return true
	})

	return
}

// Publish an event. Function will block until all subscribers are done.
// Also consider the function:
//
//	ebus.Pub(bus, ctx, 123, &myVar)
func (eb *EventBus) Pub(ctx context.Context, event Event) {
	key := eventKey{
		event: event,
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*list.AtomicList[func(context.Context)])

	if !ok {
		return
	}

	for fn := range al.Iter() {
		fn(ctx)
	}
}

// Subscribe for an event. Any publisher will be blocked until all subsribers are done,
// so please keep your subscriber fast and run anything slow in e.g. a background worker.
// Also consider the function:
//
//	ebus.Sub(bus, 123, func(ctx context.Context, myVar *myType) { ... })
func (eb *EventBus) Sub(event Event, fn func(context.Context)) func(context.Context) {
	key := eventKey{
		event: event,
	}

	v, _ := eb.cbs.LoadOrStore(key, &list.AtomicList[func(context.Context)]{})

	al, ok := v.(*list.AtomicList[func(context.Context)])

	if !ok {
		return nil
	}

	al.Add(fn)
	return fn
}

// Unsubscribe an event. Returns whether there was a subscription or not.
// Also consider the function:
//
//	ebus.Unsub(bus, 123, mySubscriber)
func (eb *EventBus) Unsub(event Event, fn func(context.Context)) (unsubscribed bool) {
	key := eventKey{
		event: event,
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*list.AtomicList[func(context.Context)])

	if !ok {
		return
	}

	return al.Remove(func(f func(context.Context)) bool {
		return same(f, fn)
	})
}
