package lru;

import (
  "container/list"
  "errors"
)

// Key for the Entity
type Key interface{}

// Value for the Entity
type Value interface{}

// Cache is lru cache
type Cache struct {
  limit   int
  items   map[Key]*list.Element
  storage *list.List
}

// Entry is lru record
type Entry struct {
  key   Key
  value Value
}

// New is creating new LRU Cache
func New(limit int) *Cache {
  if limit <= 0 {
    errors.New("Invalid limit")
    return nil
  }

  return &Cache{
    limit: limit,
    storage: list.New(),
    items: make(map[Key]*list.Element),
  }
}

// Set item to LRU
func (lru *Cache) Set(key Key, value Value) bool {
  if entry, exists := lru.items[key]; exists {
    lru.storage.MoveToFront(entry)
    lru.items[key] = entry
    return true
  }

  entry := lru.storage.PushFront(&Entry{key, value})
  lru.items[key] = entry

  if lru.overLimit() {
    lru.clanUp()
  }
  return true
}

// Check is cache over limit
func (lru *Cache) overLimit() bool {
  return lru.limit != 0 && lru.limit < lru.storage.Len()
}

// Get item from the LRU and update usage of entity
func (lru *Cache) Get(key Key) Value {
  if lru.items == nil {
    return nil
  }

  if item, ok := lru.items[key]; ok {
    lru.storage.MoveToFront(item)
    return item.Value.(*Entry).value
  }
  return nil
}

// Remove element from LRU
func (lru *Cache) remove(key Key) bool {
  if key == nil {
    return false
  }
  if el, ok := lru.items[key]; ok {
    lru.removeElement(el)
    return true
  }
  return false
}

// Remove element from list by key
func (lru *Cache) removeElement(element *list.Element) bool {
  lru.storage.Remove(element)
  if el, ok := element.Value.(*Entry); ok {
    delete(lru.items, el.key)
    return true
  }
  return false
}

// Remove eldest element in list
func (lru *Cache) clanUp() bool {
  lastElement := lru.storage.Back()

  if lastElement != nil {
    lru.removeElement(lastElement)
    return true
  }
  return false
}
