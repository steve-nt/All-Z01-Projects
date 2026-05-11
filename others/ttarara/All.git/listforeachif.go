package piscine

func IsPositiveNode(node *NodeL) bool {
	switch v := node.Data.(type) {
	case int:
		return v > 0
	case float32:
		return v > 0
	case float64:
		return v > 0
	case byte:
		return v > 0
	default:
		return false
	}
}

func IsAlNode(node *NodeL) bool {
	switch node.Data.(type) {
	case int, float32, float64, byte:
		return false
	default:
		return true
	}
}

func ListForEachIf(l *List, f func(*NodeL), cond func(*NodeL) bool) {
	for node := l.Head; node != nil; node = node.Next {
		if cond(node) {
			f(node)
		}
	}
}
