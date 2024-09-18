# ebus
A generic, thread-safe, [high-performance](#benchmark) event bus for Go. Zero allocations during pub/sub.

Please note:
- Publishers will be blocked until all subscribers of the event are done. Use [background workers](#background-workers) for any slow work.
- Any variable that is published in an event is NOT safe for a subscriber to keep after return.

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

    // Note that it's NOT safe to keep `u` after return
})

// Subscribe for created orders
ebus.Sub(bus, Created, func(o *Order) {
    fmt.Println("order", o.Id, "was created")

    // Note that it's NOT safe to keep `o` after return
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

### Background workers
As publishers are blocked until all subscribers of the event are done, a subscriber should do its work
rather fast. Any slow work should be put in the background - this is an example for how.
```go
const (
    _  ebus.Event = iota

    Created // A generic event that "something was created"
)

type Order struct{
    Id int
}

bus := ebus.NewEventBus()

// Create a buffered channel
ch := make(chan Order, 8)

go func(ch <-chan Order) {
    for {
        o, ok := <-ch

        if !ok {
            return
        }

        // Do any slow work, e.g. synchronize the order to a third-party API
        doAnySlowWork(o)
    }
}(ch)

// Subscribe for created orders and forward them to the channel
ebus.SubToChan(bus, Created, ch)

// Publish that an order was created
ebus.Pub(bus, Created, &Order{
    Id: 123456
})

// Close channel when done
close(ch)
```

## Benchmark
This gives you an idea of the performance, but your mileage may vary. [Do your own benchmark](./bench_test.go).
```
goos: darwin
goarch: arm64
pkg: github.com/webmafia/ebus
cpu: Apple M1 Pro
BenchmarkEventBus/01_subscribers-10                   65841814     17.590 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus/02_subscribers-10                   60300748     20.120 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus/04_subscribers-10                   47710713     25.010 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus/08_subscribers-10                   34022806     35.360 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus/16_subscribers-10                   16668036     72.400 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus/32_subscribers-10                   10365789    115.000 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus/64_subscribers-10                    5991422    200.300 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/01_subscribers-10        279230622      4.252 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/02_subscribers-10        259206556      4.594 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/04_subscribers-10        191008094      5.891 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/08_subscribers-10        148670781      8.141 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/16_subscribers-10         91247522     14.100 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/32_subscribers-10         62021378     19.390 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Parallell/64_subscribers-10         41622828     28.570 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/01_subscribers-10               69984446     16.860 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/02_subscribers-10               60236678     19.920 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/04_subscribers-10               45867614     26.180 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/08_subscribers-10               31973780     37.550 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/16_subscribers-10               16673738     71.790 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/32_subscribers-10               10464192    114.600 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var/64_subscribers-10                5991487    199.800 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/01_subscribers-10    262930784      4.562 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/02_subscribers-10    262294675      4.904 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/04_subscribers-10    192289729      6.090 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/08_subscribers-10    144890751      8.254 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/16_subscribers-10     90990375     13.760 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/32_subscribers-10     62902694     18.760 ns/op    0 B/op    0 allocs/op
BenchmarkEventBus_Var_Parallell/64_subscribers-10     41338268     28.810 ns/op    0 B/op    0 allocs/op
PASS
ok      github.com/webmafia/ebus       39.693s
```