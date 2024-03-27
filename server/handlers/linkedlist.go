// linkedlist.go

package handlers

import (
	// "fmt"
	"time"
)

type node struct {
	prev       *node
	next       *node
	key        string
	value      string
	ExpiryTime time.Time
}

// Define custom doubly linked list structure
type linkedList struct {
	head *node
	tail *node
}

// Add a node to the front of the linked list
func (l *linkedList) addToFront(n *node) {
	if l.head == nil {
		l.head = n
		l.tail = n
	} else {
		n.next = l.head
		l.head.prev = n
		l.head = n
	}
}

// Move a node to the front of the linked list
func (l *linkedList) moveToFront(n *node) {
	if l.head == n {
		return
	}

	if n == l.tail {
		l.tail = n.prev
	} else {
		n.next.prev = n.prev
	}
	n.prev.next = n.next

	n.next = l.head
	n.prev = nil
	l.head.prev = n
	l.head = n
}

// Remove the tail node from the linked list
func (l *linkedList) removeTail() {
	if l.tail != nil {
		if l.head == l.tail {
			l.head = nil
			l.tail = nil
		} else {
			l.tail = l.tail.prev
			l.tail.next = nil
		}
	}
}

func (c *LRUCache) removeNode(n *node) {
	if n == c.lruList.head {
		c.lruList.head = n.next
	}
	if n == c.lruList.tail {
		c.lruList.tail = n.prev
	}
	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
}
