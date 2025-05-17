package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Node struct {
	Key   int
	Value int
	Left  *Node
	Right *Node
}

type OrderedMap struct {
	root  *Node
	count int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{
		root:  nil,
		count: 0,
	}
}

func (m *OrderedMap) Insert(key, value int) {
	if m.root == nil {
		m.root = &Node{Key: key, Value: value}
		m.count++
		return
	}

	m.insertNode(m.root, key, value)
}

func (m *OrderedMap) insertNode(node *Node, key, value int) {
	if key == node.Key {
		node.Value = value
		return
	}

	if key < node.Key {
		if node.Left == nil {
			node.Left = &Node{Key: key, Value: value}
			m.count++
		} else {
			m.insertNode(node.Left, key, value)
		}
	} else {
		if node.Right == nil {
			node.Right = &Node{Key: key, Value: value}
			m.count++
		} else {
			m.insertNode(node.Right, key, value)
		}
	}
}

func (m *OrderedMap) Erase(key int) {
	if m.root == nil {
		return
	}

	m.root = m.eraseNode(m.root, key)
}

func (m *OrderedMap) eraseNode(node *Node, key int) *Node {
	if node == nil {
		return nil
	}

	if key < node.Key {
		node.Left = m.eraseNode(node.Left, key)
	} else if key > node.Key {
		node.Right = m.eraseNode(node.Right, key)
	} else {
		if node.Left == nil {
			m.count--
			return node.Right
		} else if node.Right == nil {
			m.count--
			return node.Left
		}

		successor := m.findMin(node.Right)

		node.Key = successor.Key
		node.Value = successor.Value

		node.Right = m.eraseNode(node.Right, successor.Key)
	}

	return node
}

func (m *OrderedMap) findMin(node *Node) *Node {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (m *OrderedMap) Contains(key int) bool {
	return m.findNode(m.root, key) != nil
}

func (m *OrderedMap) findNode(node *Node, key int) *Node {
	if node == nil {
		return nil
	}

	if key == node.Key {
		return node
	}

	if key < node.Key {
		return m.findNode(node.Left, key)
	}

	return m.findNode(node.Right, key)
}

func (m *OrderedMap) Size() int {
	return m.count
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.inOrderTraversal(m.root, action)
}

func (m *OrderedMap) inOrderTraversal(node *Node, action func(int, int)) {
	if node == nil {
		return
	}

	m.inOrderTraversal(node.Left, action)

	action(node.Key, node.Value)

	m.inOrderTraversal(node.Right, action)
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
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
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
