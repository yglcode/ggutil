/*
 * Wrapper types for working with std heap, sort
 */
package ggutil

import (
	"container/heap"
)

type LessFunc[T any] func(T, T) bool

// Slice wrapper type to work with "sort" pkg
type Slice[T any] struct {
	Data     []T
	LessFunc LessFunc[T]
}

func (s Slice[T]) Len() int { return len(s.Data) }
func (s Slice[T]) Less(i, j int) bool {
	return s.LessFunc(s.Data[i], s.Data[j])
}
func (s Slice[T]) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

// Heap wrapper type to work with container/heap
type Heap[T any] struct {
	Slice[T]
}

func NewHeap[T any](c LessFunc[T], data ...[]T) *Heap[T] {
	h := &Heap[T]{Slice[T]{LessFunc: c}}
	if len(data) > 0 {
		h.Data = data[0]
		heap.Init(h)
	}
	return h
}

func (h *Heap[T]) Push(x any) {
	h.Data = append(h.Data, x.(T))
}

func (h *Heap[T]) Pop() any {
	l := len(h.Data) - 1
	x := h.Data[l]
	h.Data = h.Data[:l]
	return x
}

func (h Heap[T]) Peek() T {
	return h.Data[0]
}

// HeapWithIndex works with container/heap while
// maintaining values' indices in heap, which can be
// used with heap.Fix() and heap.Remove()
type HeapWithIndex[T comparable] struct {
	Heap[T]
	//track Data index in heap
	indices map[T]int
}

func NewHeapWithIndex[T comparable](c LessFunc[T], data ...[]T) *HeapWithIndex[T] {
	h := &HeapWithIndex[T]{
		Heap[T]{Slice[T]{LessFunc: c}},
		make(map[T]int),
	}
	if len(data) > 0 {
		h.Data = data[0]
		for i := 0; i < len(h.Data); i++ {
			h.indices[h.Data[i]] = i
		}
		heap.Init(h)
	}
	return h
}

func (h *HeapWithIndex[T]) Push(x any) {
	h.Heap.Push(x)
	h.indices[x.(T)] = h.Len() - 1
}

func (h *HeapWithIndex[T]) Pop() any {
	x := h.Heap.Pop()
	delete(h.indices, x.(T))
	return x
}

func (h HeapWithIndex[T]) Swap(i, j int) {
	h.Slice.Swap(i, j)
	h.indices[h.Data[i]] = i
	h.indices[h.Data[j]] = j
}

func (h HeapWithIndex[T]) Index(x T) int {
	i, ok := h.indices[x]
	if ok {
		return i
	}
	return -1
}
