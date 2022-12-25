package hw04lrucache

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
}

/*
	type cacheItem struct {
		key   Key
		value interface{}
	}
*/

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	if _, ok := lc.items[key]; !ok {
		lc.items[key] = lc.queue.PushFront(value)
		if lc.queue.Len() > lc.capacity {
			lastElem := lc.queue.Back()
			lc.queue.Remove(lastElem)
			for k, v := range lc.items {
				if v == lastElem {
					delete(lc.items, k)
					break
				}
			}
		}
		return false
	}
	lc.items[key].Value = value
	lc.queue.MoveToFront(lc.items[key])
	return true
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := lc.items[key]; !ok {
		return nil, false
	}
	value := lc.items[key].Value
	lc.queue.MoveToFront(lc.items[key])
	return value, true
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
