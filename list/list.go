package list

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/ytz123/golist"
)

type intList struct {
	head   *intNode
	length int64
}

type intNode struct {
	value  int
	marked uint32
	next   unsafe.Pointer
	sync.Mutex
}

// isMarked 标记节点是否被删除.
func (n *intNode) isMarked() bool {
	return atomic.LoadUint32(&n.marked) != 0
}

// setMarked 标记节点未已删除.
func (n *intNode) setMarked() {
	atomic.StoreUint32(&n.marked, uint32(1))
}

// getNext 获取节点后继.
func (n *intNode) getNext() *intNode {
	return (*intNode)(atomic.LoadPointer(&n.next))
}

// setNext 设置节点后继.
func (n *intNode) setNext(next *intNode) {
	atomic.StorePointer(&n.next, unsafe.Pointer(next))
}

func newIntNode(value int) *intNode {
	return &intNode{value: value}
}

// NewInt 返回一个空的有序int链表.
func NewInt() golist.IntList {
	l := &intList{head: newIntNode(0)}
	return l
}

func (l *intList) Insert(value int) (ok bool) {
	// 0. 如果插入成功就链表元素数量+1
	defer func() {
		if ok {
			atomic.AddInt64(&l.length, 1)
		}
	}()

	for {
		// 1.1 找到相邻的ab节点，满足 a.value < value && (b == nil || value <= b.value)
		a := l.head
		b := a.getNext()
		for b != nil && b.value < value {
			a = b
			b = b.getNext()
		}

		// 1.2 b.value == value, value 已存在, 插入失败
		if b != nil && b.value == value {
			ok = false
			return
		}

		isDone := func() bool {
			// 2.1 锁a
			a.Lock()
			// 5. 解锁a
			defer a.Unlock()

			// 2.2 检查 ab相邻且a未被删除, false 则返回步骤 1
			if a.getNext() != b || a.isMarked() {
				return false
			}

			// 3. 创建新节点 c
			c := newIntNode(value)
			// 4. c.next = b; a.next = c
			c.setNext(b)
			a.setNext(c)

			return true
		}()

		if !isDone {
			continue
		}

		ok = true
		return
	}
}

func (l *intList) Delete(value int) (ok bool) {
	// 0. 如果删除成功就链表元素数量-1
	defer func() {
		if ok {
			atomic.AddInt64(&l.length, -1)
		}
	}()

	for {
		// 1.1 找到相邻的ab节点，满足 a.value < value && (b == nil || value <= b.value)
		a := l.head
		b := a.getNext()
		for b != nil && b.value < value {
			a = b
			b = b.getNext()
		}

		// 1.2 没有找到 value, 删除失败
		if b == nil || b.value != value {
			ok = false
			return
		}

		isDone := func() bool {
			// 2.1 锁b
			b.Lock()
			// 5.2 解锁b
			defer b.Unlock()

			// 2.2 检查节点 b 是否被删除, 如果已删除则返回步骤 1
			if b.isMarked() {
				return false
			}

			// 3.1 锁a
			a.Lock()
			// 5.1 解锁a
			defer a.Unlock()

			// 3.2 如果 a.next!=b or a已被删除, 则返回步骤 1
			if a.getNext() != b || a.isMarked() {
				return false
			}

			// 4. 标记 b 被删除, a.next=b.next
			b.setMarked()
			a.setNext(b.getNext())

			return true
		}()

		if !isDone {
			continue
		}

		ok = true
		return
	}
}

func (l *intList) Contains(value int) bool {
	x := l.head.getNext()
	for x != nil && x.value < value {
		x = x.getNext()
	}
	if x == nil || x.value != value {
		return false
	}
	return !x.isMarked()
}

func (l *intList) Range(f func(value int) bool) {
	x := l.head.getNext()
	for x != nil {
		if !f(x.value) {
			break
		}
		x = x.getNext()
	}
}

func (l *intList) Len() int {
	return int(atomic.LoadInt64(&l.length))
}
