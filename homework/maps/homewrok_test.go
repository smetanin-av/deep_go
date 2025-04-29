package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/constraints"
)

// go test -v homework_test.go

type mapNode[TKey constraints.Ordered, TVal any] struct {
	key         TKey
	value       TVal
	left, right *mapNode[TKey, TVal]
}

type OrderedMap[TKey constraints.Ordered, TVal any] struct {
	head *mapNode[TKey, TVal]
	size int
}

func NewOrderedMap[TKey constraints.Ordered, TVal any]() OrderedMap[TKey, TVal] {
	return OrderedMap[TKey, TVal]{}
}

func (m *OrderedMap[TKey, TVal]) Insert(key TKey, value TVal) {
	parent, node := m.findNode(key)

	if node != nil {
		node.value = value
		return
	}

	node = &mapNode[TKey, TVal]{key: key, value: value}
	if parent == nil {
		m.head = node
	} else if key < parent.key {
		parent.left = node
	} else {
		parent.right = node
	}

	m.size++
}

func (m *OrderedMap[TKey, TVal]) Erase(key TKey) {
	parent, node := m.findNode(key)
	if node == nil {
		return
	}

	var replace *mapNode[TKey, TVal]
	switch {
	case node.left == nil:
		replace = node.right
	case node.right == nil:
		replace = node.left
	case node.right.left == nil:
		replace = node.right
		replace.left = node.left
	default:
		replace = m.popMinNode(node.right)
		replace.left, replace.right = node.left, node.right
	}

	switch node {
	case m.head:
		m.head = replace
	case parent.left:
		parent.left = replace
	case parent.right:
		parent.right = replace
	}

	m.size--
}

func (m *OrderedMap[TKey, TVal]) Contains(key TKey) bool {
	_, node := m.findNode(key)
	return node != nil
}

func (m *OrderedMap[TKey, TVal]) Size() int {
	return m.size
}

func (m *OrderedMap[TKey, TVal]) ForEach(action func(TKey, TVal)) {
	if m.head != nil {
		m.traverse(m.head, action)
	}
}

func (m *OrderedMap[TKey, TVal]) findNode(key TKey) (*mapNode[TKey, TVal], *mapNode[TKey, TVal]) {
	var parent *mapNode[TKey, TVal]

	for current := m.head; current != nil; {
		if key == current.key {
			return parent, current
		}
		if key < current.key {
			parent, current = current, current.left
		} else {
			parent, current = current, current.right
		}
	}

	return parent, nil
}

func (m *OrderedMap[TKey, TVal]) popMinNode(node *mapNode[TKey, TVal]) *mapNode[TKey, TVal] {
	var parent *mapNode[TKey, TVal]

	for node.left != nil {
		parent, node = node, node.left
	}

	if parent != nil {
		parent.left = node.right
	}

	return node
}

func (m *OrderedMap[TKey, TVal]) traverse(node *mapNode[TKey, TVal], fn func(TKey, TVal)) {
	if node.left != nil {
		m.traverse(node.left, fn)
	}

	fn(node.key, node.value)

	if node.right != nil {
		m.traverse(node.right, fn)
	}
}

func TestOrderedMap(t *testing.T) {
	data := NewOrderedMap[int, int]()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})
	assert.Equal(t, []int{2, 4, 5, 10, 12, 14, 15}, keys)

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})
	assert.Equal(t, []int{4, 5, 10, 12}, keys)

	data.Erase(10)

	keys = nil
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})
	assert.Equal(t, []int{4, 5, 12}, keys)
}
