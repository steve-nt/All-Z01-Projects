package piscine

func ListPushBack(l *List, data interface{}) {
	n := &NodeL{Data: data}
	if l.Head == nil {
		l.Head = n
		return
	}
	elt := l.Head
	for elt.Next != nil {
		elt = elt.Next
	}
	elt.Next = n
}
