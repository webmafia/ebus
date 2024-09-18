package list

import (
	"iter"
	"sync/atomic"
)

// A thread-safe linked list
type AtomicList[T any] struct {
	size atomic.Int64
	link atomic.Pointer[link[T]]
}

type link[T any] struct {
	next atomic.Pointer[link[T]]
	val  T
}

// Inserts a new value at the front of the chain
func (ac *AtomicList[T]) Add(v T) {
	newLink := &link[T]{val: v}

	for {
		// Get the current head of the chain
		head := ac.link.Load()
		// Point the new link's next to the current head
		newLink.next.Store(head)

		// Atomically set the new link as the head
		if ac.link.CompareAndSwap(head, newLink) {
			break
		}
		// If CompareAndSwap fails, another goroutine updated the head, retry
	}

	ac.size.Add(1)
}

// Removes the first occurrence from the chain
func (ac *AtomicList[T]) Remove(fn func(T) bool) bool {
	var prev *link[T] // Previous link
	curr := ac.link.Load()

	for curr != nil {
		// If the current value matches the value to be removed
		if fn(curr.val) {
			next := curr.next.Load()

			if prev == nil {
				// If we are at the head, update the head to be the next link
				if ac.link.CompareAndSwap(curr, next) {
					ac.size.Add(-1)
					return true
				}
			} else {
				// Otherwise, update the previous link's next to bypass the current one
				if prev.next.CompareAndSwap(curr, next) {
					ac.size.Add(-1)
					return true
				}
			}
			return false // CAS failed, retry not needed in this simple logic
		}

		// Move to the next link
		prev = curr
		curr = curr.next.Load()
	}

	return false // Value not found
}

// Removes all occurrences from the chain
func (ac *AtomicList[T]) RemoveAll(fn func(T) bool) (removed int64) {
	var prev *link[T] // Previous link
	curr := ac.link.Load()

	for curr != nil {
		next := curr.next.Load()

		// If the current value matches the condition for removal
		if fn(curr.val) {
			if prev == nil {
				// If we're at the head, update the head to be the next link
				if ac.link.CompareAndSwap(curr, next) {
					removed++
				}
			} else {
				// Update the previous link's next to bypass the current one
				if prev.next.CompareAndSwap(curr, next) {
					removed++
				}
			}
			// Move to the next link but do NOT advance prev since it's still valid
			curr = next
		} else {
			// Move both prev and curr to the next link
			prev = curr
			curr = next
		}
	}

	ac.size.Add(-removed)
	return
}

// Returns an iterator
func (ac *AtomicList[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		l := ac.link.Load()

		for l != nil {
			if !yield(l.val) {
				break
			}

			l = l.next.Load()
		}
	}
}

// Returns the number of items in the list
func (ac *AtomicList[T]) Size() int64 {
	return ac.size.Load()
}

// Reset the list
func (ac *AtomicList[T]) Reset() {
	ac.link.Store(nil)
	ac.size.Store(0)
}
