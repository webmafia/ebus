package ebus

import (
	"context"
	"fmt"
	"testing"
)

func BenchmarkEventBus(b *testing.B) {
	const MyEvent Event = 1
	subs := []int{1, 2, 4, 8, 16, 32, 64}

	for _, n := range subs {
		b.Run(fmt.Sprintf("%02d_subscribers", n), func(b *testing.B) {
			eb := NewEventBus()
			ctx := context.Background()

			for range n {
				eb.Sub(MyEvent, func(_ context.Context) {
					_ = MyEvent
				})
			}

			b.ResetTimer()

			for range b.N {
				eb.Pub(ctx, MyEvent)
			}
		})
	}
}

func BenchmarkEventBus_Parallell(b *testing.B) {
	const MyEvent Event = 1
	subs := []int{1, 2, 4, 8, 16, 32, 64}

	for _, n := range subs {
		b.Run(fmt.Sprintf("%02d_subscribers", n), func(b *testing.B) {
			eb := NewEventBus()
			ctx := context.Background()

			for range n {
				eb.Sub(MyEvent, func(_ context.Context) {
					_ = MyEvent
				})
			}

			b.ResetTimer()

			b.RunParallel(func(p *testing.PB) {
				for p.Next() {
					eb.Pub(ctx, MyEvent)
				}
			})
		})
	}
}

func BenchmarkEventBus_Var(b *testing.B) {
	const MyEvent Event = 1
	subs := []int{1, 2, 4, 8, 16, 32, 64}

	for _, n := range subs {
		b.Run(fmt.Sprintf("%02d_subscribers", n), func(b *testing.B) {
			eb := NewEventBus()
			ctx := context.Background()

			for range n {
				Sub(eb, MyEvent, func(_ context.Context, v *int) {
					_ = v
				})
			}

			b.ResetTimer()

			for i := range b.N {
				Pub(eb, ctx, MyEvent, &i)
			}
		})
	}
}

func BenchmarkEventBus_Var_Parallell(b *testing.B) {
	const MyEvent Event = 1
	subs := []int{1, 2, 4, 8, 16, 32, 64}

	for _, n := range subs {
		b.Run(fmt.Sprintf("%02d_subscribers", n), func(b *testing.B) {
			eb := NewEventBus()
			ctx := context.Background()

			for range n {
				Sub(eb, MyEvent, func(_ context.Context, v *int) {
					_ = v
				})
			}

			b.ResetTimer()

			b.RunParallel(func(p *testing.PB) {
				var i int

				for p.Next() {
					Pub(eb, ctx, MyEvent, &i)
				}
			})
		})
	}
}
