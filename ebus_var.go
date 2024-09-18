package ebus

import "github.com/webmafia/ebus/list"

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

// Subscribe for an event with a variable. Any publisher will be blocked until all subsribers are done,
// so please keep your subscriber fast and run anything slow in e.g. a background worker.
// Subscribers should NOT keep the variable after return.
// Publishers must publish the specific variable type.
func Sub[T any](eb *EventBus, event Event, fn func(*T)) func(*T) {
	key := eventKey{
		event: event,
		typ:   typeHash((*T)(nil)),
	}

	v, _ := eb.cbs.LoadOrStore(key, &list.AtomicList[func(*T)]{})

	al, ok := v.(*list.AtomicList[func(*T)])

	if !ok {
		return nil
	}

	al.Add(fn)
	return fn
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
