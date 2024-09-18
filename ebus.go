package ebus

import (
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
//	bus.Pub(123)
//
//	// Or also send variable
//	ebus.Pub(bus, 123, &myVar)
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
//	ebus.Pub(bus, 123, &myVar)
func (eb *EventBus) Pub(event Event) {
	key := eventKey{
		event: event,
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*list.AtomicList[func()])

	if !ok {
		return
	}

	for fn := range al.Iter() {
		fn()
	}
}

// Subscribe for an event. Any publishers will be blocked until all subsribers are done,
// so please keep your subscriber fast and run anything slow in e.g. a background worker.
// Also consider the function:
//
//	ebus.Sub(bus, 123, func(myVar *myType) { ... })
func (eb *EventBus) Sub(event Event, fn func()) {
	key := eventKey{
		event: event,
	}

	v, _ := eb.cbs.LoadOrStore(key, &list.AtomicList[func()]{})

	al, ok := v.(*list.AtomicList[func()])

	if !ok {
		return
	}

	al.Add(fn)
}

// Unsubscribe an event. Returns whether there was a subscription or not.
// Also consider the function:
//
//	ebus.Unsub(bus, 123, mySubscriber)
func (eb *EventBus) Unsub(event Event, fn func()) (unsubscribed bool) {
	key := eventKey{
		event: event,
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*list.AtomicList[func()])

	if !ok {
		return
	}

	return al.Remove(func(f func()) bool {
		return same(f, fn)
	})
}

// Publish an event with a variable. Function will block until all subscribers are done.
// Subscribers must subscribe for the specific variable type.
func Pub[T any](eb *EventBus, event Event, val *T) {
	key := eventKey{
		event: event,
		typ:   typeHash(val),
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*list.AtomicList[func(*T)])

	if !ok {
		return
	}

	for fn := range al.Iter() {
		fn(noescapeVal(val))
	}
}

// Subscribe for an event. Any publishers will be blocked until all subsribers are done,
// so please keep your subscriber fast and run anything slow in e.g. a background worker.
// Publishers must publish the specific variable type.
func Sub[T any](eb *EventBus, event Event, fn func(*T)) {
	key := eventKey{
		event: event,
		typ:   typeHash((*T)(nil)),
	}

	v, _ := eb.cbs.LoadOrStore(key, &list.AtomicList[func(*T)]{})

	al, ok := v.(*list.AtomicList[func(*T)])

	if !ok {
		return
	}

	al.Add(fn)
}

// Unsubscribe an event. Returns whether there was a subscription or not.
func Unsub[T any](eb *EventBus, event Event, fn func(*T)) (unsubscribed bool) {
	key := eventKey{
		event: event,
		typ:   typeHash((*T)(nil)),
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*list.AtomicList[func(*T)])

	if !ok {
		return
	}

	return al.Remove(func(f func(*T)) bool {
		return same(f, fn)
	})
}
