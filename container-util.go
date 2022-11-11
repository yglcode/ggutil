/*
 * Wrapper types for working with std heap, sort
 */
package ggutil

import (
	"container/heap"
)

type Less[T any] func(T, T) bool

// Slice wrapper type to work with "sort" pkg
type Slice[T any] struct {
	data []T
	less Less[T]
}

func (s Slice[T]) Len() int { return len(s.data) }
func (s Slice[T]) Less(i, j int) bool {
	return s.less(s.data[i], s.data[j])
}
func (s Slice[T]) Swap(i, j int) {
	s.data[i], s.data[j] = s.data[j], s.data[i]
}

// Heap wrapper type to work with container/heap
type Heap[T any] struct {
	Slice[T]
}

func NewHeap[T any](c Less[T], data ...[]T) *Heap[T] {
	h := &Heap[T]{Slice[T]{less: c}}
	if len(data) > 0 {
		h.Slice.data = data[0]
		heap.Init(h)
	}
	return h
}

func (h *Heap[T]) Push(x interface{}) {
	h.Slice.data = append(h.Slice.data, x.(T))
}

func (h *Heap[T]) Pop() interface{} {
	l := len(h.Slice.data) - 1
	x := h.Slice.data[l]
	h.Slice.data = h.Slice.data[:l]
	return x
}

func (h Heap[T]) Peek() T {
	return h.Slice.data[0]
}

// HeapWithIndex works with container/heap while
// maintaining values' indices in heap, which can be
// used with heap.Fix() and heap.Remove()
type HeapWithIndex[T comparable] struct {
	Heap[T]
	//track data index in heap
	indices map[T]int
}

func NewHeapWithIndex[T comparable](c Less[T], data ...[]T) *HeapWithIndex[T] {
	h := &HeapWithIndex[T]{Heap[T]{Slice[T]{less: c}}, nil}
	h.indices = make(map[T]int)
	if len(data) > 0 {
		h.Slice.data = data[0]
		for i := 0; i < len(h.Slice.data); i++ {
			h.indices[h.Slice.data[i]] = i
		}
		heap.Init(h)
	}
	return h
}

func (h *HeapWithIndex[T]) Push(x interface{}) {
	h.Heap.Push(x)
	h.indices[x.(T)] = h.Len() - 1
}

func (h *HeapWithIndex[T]) Pop() interface{} {
	x := h.Heap.Pop()
	delete(h.indices, x.(T))
	return x
}

func (h HeapWithIndex[T]) Swap(i, j int) {
	h.Slice.Swap(i, j)
	h.indices[h.Slice.data[i]] = i
	h.indices[h.Slice.data[j]] = j
}

func (h HeapWithIndex[T]) Index(x T) int {
	i, ok := h.indices[x]
	if ok {
		return i
	}
	return -1
}
