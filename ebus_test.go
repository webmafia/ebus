package ebus

import (
	"fmt"
)

func ExampleEventBus() {
	const MyEvent Event = 1

	eb := NewEventBus()

	subA := func(v *int) {
		fmt.Println("Subscriber A:", *v)
	}

	subB := func(v *int) {
		fmt.Println("Subscriber B:", *v)
	}

	subC := func(v *int) {
		fmt.Println("Subscriber C:", *v)
	}

	Sub(eb, MyEvent, subA)
	Sub(eb, MyEvent, subB)
	Sub(eb, MyEvent, subC)

	fmt.Println("Subscribers:", eb.Subscribers())

	i := 123

	Pub(eb, MyEvent, &i)

	fmt.Println("--------")
	Unsub(eb, MyEvent, subB)
	fmt.Println("Subscribers:", eb.Subscribers())
	Pub(eb, MyEvent, &i)

	// Output:
	//
	// Subscribers: 3
	// Subscriber C: 123
	// Subscriber B: 123
	// Subscriber A: 123
	// --------
	// Subscribers: 2
	// Subscriber C: 123
	// Subscriber A: 123
}
