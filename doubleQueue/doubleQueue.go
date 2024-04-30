package doubleQueue

import (
	"sync"
)

type Node struct {
    value interface{}
    prev  *Node
    next  *Node
}

type DLQueue struct {
    head   *Node
    tail   *Node
    length int
    lock   sync.Mutex
}

func NewdLQueue() *DLQueue {
    return &DLQueue{}
}

// PushTop adds an item to the top of the DLQueue
func (d *DLQueue) PushTop(value interface{}) {
    d.lock.Lock()
    defer d.lock.Unlock()

    newNode := &Node{value: value}
    if d.head == nil {
        d.head, d.tail = newNode, newNode
    } else {
        newNode.next = d.head
        d.head.prev = newNode
        d.head = newNode
    }
    d.length++
}

// PopTop removes and returns an item from the top of the DLQueue
func (d *DLQueue) PopTop() interface{} {
    d.lock.Lock()
    defer d.lock.Unlock()

    if d.head == nil {
        return nil
    }

    ret := d.head.value
    d.head = d.head.next
    if d.head == nil {
        d.tail = nil
    } else {
        d.head.prev = nil
    }
    d.length--
    return ret
}

// PopBottom removes and returns an item from the bottom of the DLQueue
func (d *DLQueue) PopBottom() interface{} {
    d.lock.Lock()
    defer d.lock.Unlock()

    if d.tail == nil {
        return nil
    }

    ret := d.tail.value
    d.tail = d.tail.prev
    if d.tail == nil {
        d.head = nil
    } else {
        d.tail.next = nil
    }
    d.length--
    return ret
}

// Size returns the current size of the DLQueue
func (d *DLQueue) Size() int {
    d.lock.Lock()
    defer d.lock.Unlock()
    return d.length
}

