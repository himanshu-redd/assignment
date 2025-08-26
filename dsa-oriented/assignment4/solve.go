package main

import "fmt"

type Node struct {
	Val  int
	Next *Node
}

func reverseIterative(head *Node) *Node {
	var prev *Node
	curr := head

	for curr != nil {
		nextTemp := curr.Next 
		curr.Next = prev      
		prev = curr           
		curr = nextTemp       
	}
	return prev
}

func reverseRecursive(head *Node) *Node {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := reverseRecursive(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHead
}

func printList(head *Node) {
	for head != nil {
		fmt.Print(head.Val)
		if head.Next != nil {
			fmt.Print(" -> ")
		}
		head = head.Next
	}
	fmt.Println()
}

func main() {
	head := &Node{Val: 1}
	head.Next = &Node{Val: 2}
	head.Next.Next = &Node{Val: 3}

	fmt.Print("Original: ")
	printList(head)

	iterReversed := reverseIterative(head)
	fmt.Print("Iterative Reverse: ")
	printList(iterReversed)

	recurReversed := reverseRecursive(iterReversed)
	fmt.Print("Recursive Reverse: ")
	printList(recurReversed)
}
