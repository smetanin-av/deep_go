package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

// go test -v homework_test.go

type CircularQueue[T constraints.Signed] struct {
	values           []T
	len, front, rear int
}

func NewCircularQueue[T constraints.Signed](size int) CircularQueue[T] {
	return CircularQueue[T]{rear: -1, values: make([]T, size)}
}

func (q *CircularQueue[T]) Push(value T) bool {
	if q.Full() {
		return false
	}

	q.rear = q.nextPos(q.rear)
	q.values[q.rear] = value
	q.len++

	return true
}

func (q *CircularQueue[T]) Pop() bool {
	if q.Empty() {
		return false
	}

	q.front = q.nextPos(q.front)
	q.len--

	return true
}

func (q *CircularQueue[T]) Front() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.front]
}

func (q *CircularQueue[T]) Back() T {
	if q.Empty() {
		return -1
	}

	return q.values[q.rear]
}

func (q *CircularQueue[T]) Empty() bool {
	return q.len == 0
}

func (q *CircularQueue[T]) Full() bool {
	return q.len == q.Size()
}

func (q *CircularQueue[T]) Size() int {
	return len(q.values)
}

func (q *CircularQueue[T]) Clear() {
	q.len = 0
	q.front = 0
	q.rear = -1
}

func (q *CircularQueue[T]) Resize(newSize int) {
	if q.len > newSize {
		q.len = newSize
	}

	values := make([]T, newSize)
	for i := range q.len {
		values[i] = q.values[q.front]
		q.front = q.nextPos(q.front)
	}

	q.values = values
	q.front = 0
	q.rear = q.len - 1
}

func (q *CircularQueue[T]) nextPos(pos int) int {
	if pos == q.Size()-1 {
		return 0
	}

	return pos + 1
}

func TestCircularQueue(t *testing.T) {
	t.Run("critical path", func(t *testing.T) {
		const queueSize = 3
		queue := NewCircularQueue[int](queueSize)

		assert.True(t, queue.Empty())
		assert.False(t, queue.Full())

		assert.Equal(t, -1, queue.Front())
		assert.Equal(t, -1, queue.Back())
		assert.False(t, queue.Pop())

		assert.True(t, queue.Push(1))
		assert.True(t, queue.Push(2))
		assert.True(t, queue.Push(3))
		assert.False(t, queue.Push(4))

		assert.True(t, reflect.DeepEqual([]int{1, 2, 3}, queue.values))

		assert.False(t, queue.Empty())
		assert.True(t, queue.Full())

		assert.Equal(t, 1, queue.Front())
		assert.Equal(t, 3, queue.Back())

		assert.True(t, queue.Pop())
		assert.False(t, queue.Empty())
		assert.False(t, queue.Full())
		assert.True(t, queue.Push(4))

		assert.True(t, reflect.DeepEqual([]int{4, 2, 3}, queue.values))

		assert.Equal(t, 2, queue.Front())
		assert.Equal(t, 4, queue.Back())

		assert.True(t, queue.Pop())
		assert.True(t, queue.Pop())
		assert.True(t, queue.Pop())
		assert.False(t, queue.Pop())

		assert.True(t, queue.Empty())
		assert.False(t, queue.Full())
	})

	t.Run("clear queue", func(t *testing.T) {
		const queueSize = 3
		queue := NewCircularQueue[int16](queueSize)

		assert.True(t, queue.Push(1))
		assert.True(t, queue.Push(2))
		assert.True(t, queue.Push(3))
		assert.Equal(t, []int16{1, 2, 3}, queue.values)

		queue.Clear()

		assert.True(t, queue.Push(4))
		assert.True(t, queue.Push(5))
		assert.True(t, queue.Push(6))
		assert.Equal(t, []int16{4, 5, 6}, queue.values)
		assert.Equal(t, int16(4), queue.Front())
		assert.Equal(t, int16(6), queue.Back())
	})

	t.Run("resize queue", func(t *testing.T) {
		const initQueueSize = 2
		queue := NewCircularQueue[int32](initQueueSize)
		assert.Equal(t, initQueueSize, queue.Size())

		assert.True(t, queue.Push(1))
		assert.True(t, queue.Push(2))
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(3))
		assert.True(t, queue.Full())
		assert.Equal(t, []int32{3, 2}, queue.values)
		assert.Equal(t, int32(2), queue.Front())
		assert.Equal(t, int32(3), queue.Back())

		const newQueueSize = initQueueSize + 2
		queue.Resize(newQueueSize)
		assert.Equal(t, newQueueSize, queue.Size())
		assert.False(t, queue.Full())
		assert.Equal(t, []int32{2, 3, 0, 0}, queue.values)
		assert.Equal(t, int32(2), queue.Front())
		assert.Equal(t, int32(3), queue.Back())

		assert.True(t, queue.Push(4))
		assert.True(t, queue.Push(5))
		assert.True(t, queue.Pop())
		assert.True(t, queue.Push(6))
		assert.True(t, queue.Full())
		assert.Equal(t, []int32{6, 3, 4, 5}, queue.values)
		assert.Equal(t, int32(3), queue.Front())
		assert.Equal(t, int32(6), queue.Back())

		queue.Resize(initQueueSize)
		assert.Equal(t, initQueueSize, queue.Size())
		assert.True(t, queue.Full())
		assert.Equal(t, []int32{3, 4}, queue.values)
		assert.Equal(t, int32(3), queue.Front())
		assert.Equal(t, int32(4), queue.Back())
	})
}
