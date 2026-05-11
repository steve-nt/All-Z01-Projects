package piscine

// TreeNode represents a node in the binary search tree.
type TreeNode struct {
	Left, Right, Parent *TreeNode
	Data                string
}

// BTreeIsBinary checks if the tree is a valid binary search tree.
func BTreeIsBinary(root *TreeNode) bool {
	return isBSTUtil(root, nil, nil)
}

// isBSTUtil is a utility function to check if the tree is a valid BST.
func isBSTUtil(node, min, max *TreeNode) bool {
	// An empty tree is a BST
	if node == nil {
		return true
	}

	// If this node violates the min/max constraint, return false
	if (min != nil && node.Data <= min.Data) || (max != nil && node.Data >= max.Data) {
		return false
	}

	// Otherwise, check the subtrees recursively, tightening the min or max constraint
	return isBSTUtil(node.Left, min, node) && isBSTUtil(node.Right, node, max)
}

// BTreeDeleteNode deletes the node from the tree given by root.
func BTreeDeleteNode(root, node *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}

	// Find the node to delete
	if node.Data < root.Data {
		root.Left = BTreeDeleteNode(root.Left, node)
	} else if node.Data > root.Data {
		root.Right = BTreeDeleteNode(root.Right, node)
	} else {
		// Node found
		if root.Left == nil {
			temp := root.Right
			root = nil
			return temp
		} else if root.Right == nil {
			temp := root.Left
			root = nil
			return temp
		}

		// Node with two children: Get the in-order successor (smallest in the right subtree)
		temp := minValueNode(root.Right)
		root.Data = temp.Data
		root.Right = BTreeDeleteNode(root.Right, temp)
	}
	return root
}

// minValueNode finds the node with the minimum value in the tree.
func minValueNode(node *TreeNode) *TreeNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

// BTreeInsertData inserts a new node with the given data into the tree.
func BTreeInsertData(root *TreeNode, data string) *TreeNode {
	if root == nil {
		return &TreeNode{Data: data}
	}
	if data < root.Data {
		left := BTreeInsertData(root.Left, data)
		root.Left = left
		left.Parent = root
	} else {
		right := BTreeInsertData(root.Right, data)
		root.Right = right
		right.Parent = root
	}
	return root
}

// BTreeSearchItem searches for a node with the given data in the tree.
func BTreeSearchItem(root *TreeNode, data string) *TreeNode {
	if root == nil || root.Data == data {
		return root
	}
	if data < root.Data {
		return BTreeSearchItem(root.Left, data)
	}
	return BTreeSearchItem(root.Right, data)
}

// BTreeApplyInorder applies a function to each node of the tree in-order.
func BTreeApplyInorder(root *TreeNode, f func(item string)) {
	if root == nil {
		return
	}
	BTreeApplyInorder(root.Left, f)
	f(root.Data)
	BTreeApplyInorder(root.Right, f)
}
