package piscine

// ListLast returns the Data interface{} of the last element in the linked list.
// If the list is empty, it returns nil.
func ListLast(l *List) interface{} {
	if l == nil || l.Head == nil {
		return nil
	}
	// Efficient approach using the Tail pointer:
	return l.Tail.Data
}
