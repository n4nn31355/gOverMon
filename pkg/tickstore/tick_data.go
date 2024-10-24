package tickstore

import (
	"fmt"
	"iter"
	"slices"
)

type TickData[T any] struct {
	values []T
}

func NewTickData[T any](size int) *TickData[T] {
	if size < 1 {
		panic("size must be greater than zero")
	}
	return &TickData[T]{
		values: make([]T, size),
	}
}

func (td *TickData[T]) Iter() iter.Seq2[int, T] {
	return slices.All(td.values)
}

func (td *TickData[T]) IterBackward() iter.Seq2[int, T] {
	return slices.Backward(td.values)
}

func (td *TickData[T]) GetValues() []T {
	return td.values
}

func (td *TickData[T]) GetValue(idx int) T {
	return td.values[idx]
}

// TODO: First value actually is Last added? Confusing?
func (td *TickData[T]) GetFirstValue() T {
	return td.values[0]
}

func (td *TickData[T]) GetLastValue() T {
	return td.values[len(td.values)-1]
}

func (td *TickData[T]) GetSize() int {
	return len(td.values)
}

// TODO: Implement circular array for better perf(avoid allocations)?
func (td *TickData[T]) AddValues(values ...T) error {
	vLen := len(values)
	if cap(td.values) < vLen {
		return fmt.Errorf(
			"attempt to add %v values to array with size: %v",
			vLen, cap(td.values),
		)
	}

	if vLen != len(td.values) {
		td.values = slices.Replace(
			td.values, vLen, len(td.values), td.values[:len(td.values)-vLen]...)
	}
	td.values = slices.Replace(td.values, 0, vLen, values...)

	return nil
}
