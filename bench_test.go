package ebus

import (
	"sync"
	"sync/atomic"
	"testing"
)

type mutexMap struct {
	data map[int]func()
	mu   sync.RWMutex
}

func BenchmarkMutexMap(b *testing.B) {
	m := mutexMap{
		data: make(map[int]func()),
	}
	f := func() {}

	b.Run("WriteParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				m.mu.Lock()
				m.data[i] = f
				m.mu.Unlock()

				i++
			}
		})
	})

	b.Run("ReadParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				m.mu.RLock()
				v, ok := m.data[i]
				m.mu.RUnlock()

				_, _ = v, ok
				i++
			}
		})
	})
}

func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	f := func() {}

	b.Run("WriteParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				m.Store(i, f)

				i++
			}
		})
	})

	b.Run("ReadParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				v, ok := m.Load(i)

				_, _ = v, ok

				v2, ok := v.(func())

				_, _ = v2, ok
				i++
			}
		})
	})
}

func BenchmarkMap(b *testing.B) {
	m := make(map[int]func())
	f := func() {}

	b.Run("Write", func(b *testing.B) {
		for i := range b.N {
			m[i] = f
		}
	})

	b.Run("ReadParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				v, ok := m[i]

				_, _ = v, ok
				i++
			}
		})
	})
}

func BenchmarkAtomicPointer(b *testing.B) {
	var ptr atomic.Pointer[int]
	var v int

	b.Run("WriteParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				ptr.Store(&v)

				i++
			}
		})
	})

	b.Run("ReadParallell", func(b *testing.B) {
		b.RunParallel(func(p *testing.PB) {
			var i int

			for p.Next() {
				v := ptr.Load()
				_ = v
				i++
			}
		})
	})
}
