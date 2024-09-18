package ebus

import (
	"sync"
	"sync/atomic"
)

type Event uint32

type EventBus struct {
	cbs sync.Map
}

type eventKey struct {
	event Event
	typ   uint32
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

func Pub[T any](eb *EventBus, event Event, val *T) {
	key := eventKey{
		event: event,
		typ:   typeHash(val),
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*AtomicList[func(*T)])

	if !ok {
		return
	}

	for fn := range al.Iter() {
		fn(noescapeVal(val))
	}
}

func Sub[T any](eb *EventBus, event Event, fn func(*T)) {
	key := eventKey{
		event: event,
		typ:   typeHash((*T)(nil)),
	}

	v, _ := eb.cbs.LoadOrStore(key, &AtomicList[func(*T)]{})

	al, ok := v.(*AtomicList[func(*T)])

	if !ok {
		return
	}

	al.Add(fn)
}

func Unsub[T any](eb *EventBus, event Event, fn func(*T)) (unsubscribed bool) {
	key := eventKey{
		event: event,
		typ:   typeHash((*T)(nil)),
	}

	v, ok := eb.cbs.Load(key)

	if !ok {
		return
	}

	al, ok := v.(*AtomicList[func(*T)])

	if !ok {
		return
	}

	return al.Remove(func(f func(*T)) bool {
		return same(f, fn)
	})
}

func (eb *EventBus) Subscribers() (n int64) {
	eb.cbs.Range(func(_, val any) bool {
		ptr := (*atomic.Int64)(efaceData(val))
		n += ptr.Load()
		return true
	})

	return
}
