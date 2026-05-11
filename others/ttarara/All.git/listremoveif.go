package piscine

func ListRemoveIf(l *List, data_ref interface{}) {
	// Handle the case where the list is empty
	if l.Head == nil {
		return
	}

	// Remove any matching nodes at the beginning of the list
	for l.Head != nil && l.Head.Data == data_ref {
		l.Head = l.Head.Next
	}

	// If the list is now empty, update the tail and return
	if l.Head == nil {
		l.Tail = nil
		return
	}

	// Now remove matching nodes from the rest of the list
	prev := l.Head
	for curr := l.Head.Next; curr != nil; curr = curr.Next {
		if curr.Data == data_ref {
			prev.Next = curr.Next
			if curr.Next == nil { // If we removed the tail node, update the tail
				l.Tail = prev
			}
		} else {
			prev = curr
		}
	}
}
