// Package list implements a doubly-linked list. The API and behaviour match
// the container/list package in the Go version 1 standard library.
// For per function and per type documentation, see that package.
package list

type Element[T any] struct {
	prev  *Element[T]
	next  *Element[T]
	list  *List[T]
	Value T
}

func (e *Element[T]) Next() *Element[T] {
	if e.list == nil {
		return nil
	}
	if e.next == &e.list.root {
		return nil
	}
	return e.next
}

func (e *Element[T]) Prev() *Element[T] {
	if e.list == nil {
		return nil
	}
	if e.prev == &e.list.root {
		return nil
	}
	return e.prev
}

type List[T any] struct {
	root Element[T]
	size int
}

func New[T any]() *List[T] {
	return new(List[T]).Init()
}

func (l *List[T]) Back() *Element[T] {
	if l.size == 0 {
		return nil
	}
	return l.root.prev
}

func (l *List[T]) Front() *Element[T] {
	if l.size == 0 {
		return nil
	}
	return l.root.next
}

func (l *List[T]) Init() *List[T] {
	l.root.prev = &l.root
	l.root.next = &l.root
	l.size = 0
	return l
}

func (l *List[T]) insertValueAfter(v T, mark *Element[T]) *Element[T] {
	e := Element[T]{prev: mark, next: mark.next, Value: v, list: l}
	mark.next.prev = &e
	mark.next = &e
	l.size++
	return &e
}

func (l *List[T]) InsertAfter(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValueAfter(v, mark)
}

func (l *List[T]) InsertBefore(v T, mark *Element[T]) *Element[T] {
	if mark.list != l {
		return nil
	}
	return l.insertValueAfter(v, mark.prev)
}

func (l *List[T]) Len() int {
	return l.size
}

func (l *List[T]) moveAfter(e, mark *Element[T]) {
	if e == mark {
		return
	}
	// fixup around old |e| position.
	e.prev.next = e.next
	e.next.prev = e.prev
	// write new relatives for |e|.
	e.prev = mark
	e.next = mark.next
	// fixup around new |e| position.
	mark.next.prev = e
	mark.next = e
}

func (l *List[T]) MoveAfter(e, mark *Element[T]) {
	if e.list != l || mark.list != l || e == mark {
		return
	}
	l.moveAfter(e, mark)
}

func (l *List[T]) MoveBefore(e, mark *Element[T]) {
	if e.list != l || mark.list != l || e == mark {
		return
	}
	l.moveAfter(e, mark.prev)
}

func (l *List[T]) MoveToBack(e *Element[T]) {
	if e.list != l || l.root.prev == e {
		return
	}
	l.moveAfter(e, l.root.prev)
}

func (l *List[T]) MoveToFront(e *Element[T]) {
	if e.list != l || l.root.next == e {
		return
	}
	l.moveAfter(e, &l.root)
}

func (l *List[T]) PushBack(v T) *Element[T] {
	l.lazyInit()
	return l.insertValueAfter(v, l.root.prev)
}

func (l *List[T]) PushBackList(other *List[T]) {
	l.lazyInit()

	// important to "capture" other.Len() first
	// since it may be modified during the loop
	// if l == other (which is allowed).

	e := other.Front()
	for i := other.Len(); i > 0; i-- {
		l.insertValueAfter(e.Value, l.root.prev)
		e = e.Next()
	}
}

func (l *List[T]) PushFront(v T) *Element[T] {
	l.lazyInit()
	return l.insertValueAfter(v, &l.root)
}

func (l *List[T]) PushFrontList(other *List[T]) {
	l.lazyInit()

	e := other.Back()
	for i := other.Len(); i > 0; i-- {
		l.insertValueAfter(e.Value, &l.root)
		e = e.Prev()
	}
}

func (l *List[T]) Remove(e *Element[T]) T {
	if e.list != l {
		return e.Value
	}
	e.prev.next = e.next
	e.next.prev = e.prev
	e.prev = nil
	e.next = nil
	e.list = nil
	l.size--
	return e.Value
}

// lazyInit lazily initializes a zero List value.
func (l *List[T]) lazyInit() {
	if l.root.next == nil {
		l.Init()
	}
}
