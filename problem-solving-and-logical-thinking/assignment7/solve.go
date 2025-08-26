package main

import (
    "fmt"
)

type ListNode struct {
    Val  int
    Next *ListNode
}

// Merges two sorted linked lists.
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    current := dummy

    for l1 != nil && l2 != nil {
        if l1.Val < l2.Val {
            current.Next = l1
            l1 = l1.Next
        } else {
            current.Next = l2
            l2 = l2.Next
        }
        current = current.Next
    }

    if l1 != nil {
        current.Next = l1
    }
    if l2 != nil {
        current.Next = l2
    }

    return dummy.Next
}

// Merges K sorted linked lists using divide and conquer.
func mergeKLists(lists []*ListNode) *ListNode {
    if len(lists) == 0 {
        return nil
    }

    interval := 1
    for interval < len(lists) {
        for i := 0; i+interval < len(lists); i += interval * 2 {
            lists[i] = mergeTwoLists(lists[i], lists[i+interval])
        }
        interval *= 2
    }

    return lists[0]
}

func printList(node *ListNode) {
    for node != nil {
        fmt.Printf("%d -> ", node.Val)
        node = node.Next
    }
    fmt.Println("nil")
}

func main() {
    list1 := &ListNode{Val: 1, Next: &ListNode{Val: 4, Next: &ListNode{Val: 5}}}
    list2 := &ListNode{Val: 1, Next: &ListNode{Val: 3, Next: &ListNode{Val: 4}}}
    list3 := &ListNode{Val: 2, Next: &ListNode{Val: 6}}

    lists := []*ListNode{list1, list2, list3}

    mergedList := mergeKLists(lists)
    printList(mergedList)
}