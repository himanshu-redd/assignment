package main

import "fmt"

type Node struct {
	key, value int
	prev, next *Node
}

type LRUCache struct {
	capacity int
	cache    map[int]*Node
	head     *Node
	tail     *Node
}

func newLRUCache(capacity int) *LRUCache {
	h, t := &Node{}, &Node{}
	h.next = t
	t.prev = h
	return &LRUCache{capacity: capacity, cache: make(map[int]*Node), head: h, tail: t}
}

func (this *LRUCache) remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

func (this *LRUCache) insert(node *Node) {
	node.next = this.head.next
	node.prev = this.head
	this.head.next.prev = node
	this.head.next = node
}

func (this *LRUCache) Get(key int) int {
	if n, ok := this.cache[key]; ok {
		this.remove(n)
		this.insert(n)
		return n.value
	}
	return -1
}

func (this *LRUCache) Put(key int, value int) {
	if n, ok := this.cache[key]; ok {
		this.remove(n)
		delete(this.cache, key)
	}
	if len(this.cache) == this.capacity {
		lru := this.tail.prev
		this.remove(lru)
		delete(this.cache, lru.key)
	}
	n := &Node{key: key, value: value}
	this.insert(n)
	this.cache[key] = n
}

func main() {
	lru := newLRUCache(2)
	lru.Put(1, 1)
	lru.Put(2, 2)
	fmt.Println(lru.Get(1))
	lru.Put(3, 3)
	fmt.Println(lru.Get(2))
}
