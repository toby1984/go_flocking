package list

type node[T interface{}] struct {
	next  *node[T]
	prev  *node[T]
	value T
}

type List[T interface{}] struct {
	head     *node[T]
	tail     *node[T]
	size     int32
	Equality func(T, T) bool
}

func (l *List[T]) Add(value T) {
	var newNode node[T]
	newNode.value = value
	if l.head == nil {
		l.head = &newNode
		l.tail = &newNode
	} else {
		oldTail := l.tail
		oldTail.next = &newNode
		newNode.prev = oldTail
		l.tail = &newNode
	}
	l.size++
}

func (l *List[T]) ToSlice() []T {
	var result []T
	it := newListIter(l)
	for it.HasNext() {
		result = append(result, it.Next())
	}
	return result
}

func (l *List[T]) IsEmpty() bool {
	return l.Size() == 0
}

func (l *List[T]) Size() int32 {
	return l.size
}

type listIter[T interface{}] struct {
	list              *List[T]
	previouslyVisited *node[T]
	currentNode       *node[T]
}

func (it *listIter[T]) HasNext() bool {
	return it.currentNode != nil
}

func (it *listIter[T]) Remove() {
	if it.previouslyVisited == nil {
		panic("Cannot Remove() without calling Next() at least once")
	}
	it.list.removeNode(it.previouslyVisited)
}

func (it *listIter[T]) Next() T {
	if it.currentNode == nil {
		panic("Iterated past end of list")
	}
	result := it.currentNode
	it.previouslyVisited = result
	it.currentNode = it.currentNode.next
	return result.value
}

func newListIter[T interface{}](l *List[T]) listIter[T] {
	var head *node[T]
	if l.IsEmpty() {
		head = nil
	} else {
		head = l.head
	}
	return listIter[T]{list: l, currentNode: head, previouslyVisited: nil}
}

func (l *List[T]) removeNode(n *node[T]) {

	if l.IsEmpty() {
		panic("Node is not part of this list")
	}

	if l.size == 1 {
		if n != l.head {
			panic("Node is not part of this list")
		}
		l.head = nil
		l.tail = nil
	} else {
		// at least two elements
		if l.head == n {
			// remove head
			succ := n.next
			l.head = succ
			succ.prev = nil
		} else if l.tail == n {
			// remove tail
			pred := n.prev
			l.tail = pred
			pred.next = nil
			pred.prev = pred.prev
		} else {
			// remove node in-between
			pred := n.prev
			succ := n.next
			succ.prev = pred
			pred.next = succ
		}
	}

	l.size--
	n.next = nil
	n.prev = nil
}

func (l *List[T]) Contains(value T) bool {
	if l.Equality == nil {
		panic("No equality function assigned")
	}
	iter := newListIter(l)
	for iter.HasNext() {
		if l.Equality(iter.Next(), value) {
			return true
		}
	}
	return false
}

func (l *List[T]) RemoveValue(value T) bool {
	if l.Equality == nil {
		panic("No equality function assigned")
	}
	return l.Remove(func(v T) bool {
		return l.Equality(v, value)
	})
}

func (l *List[T]) Visit(visitor func(T)) {
	iter := newListIter(l)
	for iter.HasNext() {
		visitor(iter.Next())
	}
}

func (l *List[T]) Remove(predicate func(T) bool) bool {
	iter := newListIter(l)
	removed := false
	for iter.HasNext() {
		v := iter.Next()
		if predicate(v) {
			removed = true
			iter.Remove()
		}
	}
	return removed
}
