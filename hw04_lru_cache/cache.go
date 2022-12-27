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

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	var item cacheItem
	var removeItem cacheItem
	item.key = key
	item.value = value
	if _, ok := lc.items[key]; !ok {
		lc.items[key] = lc.queue.PushFront(item)
		if lc.queue.Len() > lc.capacity {
			lastElem := lc.queue.Back()
			lc.queue.Remove(lastElem)
			removeItem = lastElem.Value.(cacheItem)
			delete(lc.items, removeItem.key)
		}
		return false
	}
	lc.items[key].Value = item
	lc.queue.MoveToFront(lc.items[key])
	return true
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	var item cacheItem
	v, ok := lc.items[key]
	if !ok {
		return nil, false
	}
	item = v.Value.(cacheItem)
	lc.queue.MoveToFront(lc.items[key])
	return item.value, true
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
