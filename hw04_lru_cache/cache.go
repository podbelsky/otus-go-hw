package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mux      *sync.Mutex
}

type pair struct {
	v interface{}
	k Key
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if l.capacity == 0 {
		return false
	}

	l.mux.Lock()
	defer l.mux.Unlock()

	if item, exists := l.items[key]; exists {
		item.Value.(*pair).v = value
		_ = l.queue.MoveToFront(item) // skip error

		return true
	}

	item := l.queue.PushFront(&pair{v: value, k: key})
	l.items[key] = item

	// remove last el if enrich capacity
	if l.queue.Len() > l.capacity {
		last := l.queue.Back()
		delete(l.items, last.Value.(*pair).k)
		_ = l.queue.Remove(last) // skip error
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if item, exists := l.items[key]; exists {
		_ = l.queue.MoveToFront(item) // skip error

		return item.Value.(*pair).v, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	if l.queue.Len() == 0 {
		return
	}

	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
		mux:      &sync.Mutex{},
	}
}
