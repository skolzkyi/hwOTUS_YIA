package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	length int
	head   *ListItem
	tail   *ListItem
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	elem := ListItem{}
	elem.Value = v
	if l.length == 0 {
		l.tail = &elem
	} else {
		l.Front().Prev = &elem
		elem.Next = l.Front()
	}
	l.head = &elem
	l.length++

	return &elem
}

func (l *list) PushBack(v interface{}) *ListItem {
	elem := ListItem{}
	elem.Value = v
	if l.length == 0 {
		l.head = &elem
	} else {
		l.Back().Next = &elem
		elem.Prev = l.Back()
	}
	l.tail = &elem
	l.length++

	return &elem
}

func (l *list) Remove(i *ListItem) {
	if i.Next == nil {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	if i.Prev == nil {
		l.head = i.Next
	} else {
		i.Prev.Next = i.Next
	}
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		return
	}
	if i.Next == nil {
		l.tail = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.head
	l.head = i
}
