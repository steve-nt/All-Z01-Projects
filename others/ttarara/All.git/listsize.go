package piscine

type NodeL struct {
	Data interface{}
	Next *NodeL
}

type List struct {
	Head *NodeL
	Tail *NodeL
}

func ListSize(l *List) int {
	n := l.Head
	size := 0
	for n != nil {
		size++
		n = n.Next
	}
	return size
}
