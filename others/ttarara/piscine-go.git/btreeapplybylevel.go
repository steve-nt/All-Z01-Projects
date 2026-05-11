package piscine

// BTreeApplyByLevel applies the function f to each node of the tree given by root in level-order.
func BTreeApplyByLevel(root *TreeNode, f func(...interface{}) (int, error)) {
	if root == nil {
		return
	}

	// Initialize the queue with the root node
	queue := []*TreeNode{root}

	for len(queue) > 0 {
		// Get the first node from the queue
		current := queue[0]
		queue = queue[1:]

		// Apply the function f to the current node's data
		_, err := f(current.Data)
		if err != nil {
			return
		}

		// Enqueue the left child if it exists
		if current.Left != nil {
			queue = append(queue, current.Left)
		}

		// Enqueue the right child if it exists
		if current.Right != nil {
			queue = append(queue, current.Right)
		}
	}
}
