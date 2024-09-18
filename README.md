# ebus
A generic, thread-safe, highly optimized event bus for Go.

Please note that **publishers will block**, so you might want to put any long-running code in a 

## Install
```
go get github.com/webmafia/ebus
```

## Usage

### Only events
When just an anonymous event is needed.
```go
const (
	_  ebus.Event = iota

	MyEvent
)

bus := ebus.NewEventBus()

// Subscribe for MyEvent
sub := bus.Sub(MyEvent, func() {
	fmt.Println("recieved event")
})

// Publish MyEvent
bus.Pub(MyEvent)

// Unsubscribe for MyEvent
bus.Unsub(MyEvent, sub)
```

### Events + variable types
This is a powerful way to combine events and variable types
```go
const (
	_  ebus.Event = iota

	Created // A generic event that "something was created"
)

// These types usually belongs to the domain
type (
	User struct{
		Name string
	}

	Order struct{
		Id int
	}
)

bus := ebus.NewEventBus()

// Subscribe for created users
ebus.Sub(bus, Created, func(u *User) {
	fmt.Println("user", u.Name, "was created")
})

// Subscribe for created orders
ebus.Sub(bus, Created, func(o *Order) {
	fmt.Println("order", o.Id, "was created")
})

// Publish that a user was created
ebus.Pub(bus, Created, &User{
	Name: "John Doe"
})

// Publish that an order was created
ebus.Pub(bus, Created, &Order{
	Id: 123456
})
```

### Slow subscribers: Use channel
As publishers are blocked until the 