package main

import (
	"fmt"
)

type TreeNode struct {
	val   int
	left  *TreeNode
	right *TreeNode
	parent *TreeNode
}

type MinHeap struct {
	root  *TreeNode
	nodes []*TreeNode 
}

func (h *MinHeap) Insert(val int) {
	newNode := &TreeNode{val: val}
	if h.root == nil {
		h.root = newNode
		h.nodes = append(h.nodes, newNode)
		return
	}

	parent := h.nodes[(len(h.nodes)-1)/2]
	newNode.parent = parent
	if parent.left == nil {
		parent.left = newNode
	} else {
		parent.right = newNode
	}
	h.nodes = append(h.nodes, newNode)

	h.heapifyUp(newNode)
}

func (h *MinHeap) ExtractMin() (int, bool) {
	if h.root == nil {
		return 0, false
	}
	minVal := h.root.val

	last := h.nodes[len(h.nodes)-1]
	h.nodes = h.nodes[:len(h.nodes)-1]

	if last == h.root {
		h.root = nil
		return minVal, true
	}

	h.root.val = last.val
	if last.parent.left == last {
		last.parent.left = nil
	} else {
		last.parent.right = nil
	}

	h.heapifyDown(h.root)

	return minVal, true
}

func (h *MinHeap) heapifyUp(node *TreeNode) {
	for node.parent != nil && node.val < node.parent.val {
		node.val, node.parent.val = node.parent.val, node.val
		node = node.parent
	}
}

func (h *MinHeap) heapifyDown(node *TreeNode) {
	for node != nil {
		smallest := node
		if node.left != nil && node.left.val < smallest.val {
			smallest = node.left
		}
		if node.right != nil && node.right.val < smallest.val {
			smallest = node.right
		}
		if smallest != node {
			node.val, smallest.val = smallest.val, node.val
			node = smallest
		} else {
			break
		}
	}
}

func main() {
	h := &MinHeap{}
	h.Insert(5)
	h.Insert(3)
	h.Insert(8)

	fmt.Println(h.ExtractMin()) // 3 true
	fmt.Println(h.ExtractMin()) // 5 true
	fmt.Println(h.ExtractMin()) // 8 true
	fmt.Println(h.ExtractMin()) // 0 false
}
