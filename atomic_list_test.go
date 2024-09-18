package ebus

import (
	"fmt"
	"testing"
)

func ExampleAtomicList() {
	var ac AtomicList[int]

	for i := range 10 {
		ac.Add(i)
	}

	fmt.Println("Size:", ac.Size())

	for v := range ac.Iter() {
		fmt.Println(v)
	}

	fmt.Println("--------")

	fmt.Println(ac.Remove(func(i int) bool { return i == 4 }))
	fmt.Println(ac.Remove(func(i int) bool { return i == 9 }))
	fmt.Println(ac.Remove(func(i int) bool { return i == 0 }))
	fmt.Println(ac.Remove(func(i int) bool { return i == 0 }))
	fmt.Println("Size:", ac.Size())

	fmt.Println("--------")

	for v := range ac.Iter() {
		fmt.Println(v)
	}

	fmt.Println("--------")

	fmt.Println(ac.RemoveAll(func(i int) bool { return i < 6 }))
	fmt.Println("Size:", ac.Size())

	fmt.Println("--------")

	for v := range ac.Iter() {
		fmt.Println(v)
	}

	fmt.Println("--------")

	fmt.Println(ac.RemoveAll(func(i int) bool { return true }))
	fmt.Println("Size:", ac.Size())

	// Output:
	//
	// Size: 10
	// 9
	// 8
	// 7
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
	// 0
	// --------
	// true
	// true
	// true
	// false
	// Size: 7
	// --------
	// 8
	// 7
	// 6
	// 5
	// 3
	// 2
	// 1
	// --------
	// 4
	// Size: 3
	// --------
	// 8
	// 7
	// 6
	// --------
	// 3
	// Size: 0
}

func BenchmarkAtomicList(b *testing.B) {
	var al AtomicList[int]

	b.Run("Add1", func(b *testing.B) {
		for i := range b.N {
			al.Add(i)
		}
	})

	al.Reset()

	b.Run("Add8", func(b *testing.B) {
		for range b.N {
			for i := range 8 {
				al.Add(i)
			}
		}
	})

	al.Reset()

	b.Run("Add16", func(b *testing.B) {
		for range b.N {
			for i := range 16 {
				al.Add(i)
			}
		}
	})

	al.Reset()

	b.Run("Iter16", func(b *testing.B) {
		for i := range 16 {
			al.Add(i)
		}

		b.ResetTimer()

		for range b.N {
			for i := range al.Iter() {
				_ = i
			}
		}
	})
}

func BenchmarkAtomicListParallell(b *testing.B) {
	var al AtomicList[int]

	for i := range 16 {
		al.Add(i)
	}

	b.ResetTimer()

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			for i := range al.Iter() {
				_ = i
			}
		}
	})
}
