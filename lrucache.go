package module13__06_03

import (
	"container/list"
	"log"
	"sync"
)

type Item struct {
	Key   string
	Value interface{}
}

type LRUCache interface {
	// Add Добавляет новое значение с ключом в кеш (с наивысшим приоритетом), возвращает true, если все прошло успешно
	// В случае дублирования ключа вернуть false
	// В случае превышения размера - вытесняется наименее приоритетный элемент
	Add(key string, value interface{}) bool

	// Get Возвращает значение под ключом и флаг его наличия в кеше
	// В случае наличия в кеше элемента повышает его приоритет
	Get(key string) (value interface{}, ok bool)

	// Remove Удаляет элемент из кеша, в случае успеха возвращает true, в случае отсутствия элемента - false
	Remove(key string) (ok bool)

	getItems() map[string]*list.Element

	getPriorityQueue() *list.List
}

type LRUcache struct {
	capacity      int
	items         map[string]*list.Element
	lock          sync.RWMutex
	priorityQueue *list.List
}

func NewLRUCache(n int) LRUCache {
	if n <= 0 {
		log.Fatal("capacity should be qreater than 0")
		return nil
	}

	return &LRUcache{
		capacity:      n,
		items:         make(map[string]*list.Element, 0),
		priorityQueue: list.New(),
	}
}

func (c *LRUcache) Add(key string, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	if el, exists := c.items[key]; exists == true {
		c.priorityQueue.MoveToFront(el)
		el.Value.(*Item).Value = value
		return true
	}

	if c.priorityQueue.Len() == c.capacity {
		item := c.priorityQueue.Remove(c.priorityQueue.Back()).(*Item)
		delete(c.items, item.Key)
	}

	item := &Item{
		Key:   key,
		Value: value,
	}

	el := c.priorityQueue.PushFront(item)
	c.items[item.Key] = el

	return true
}

func (c *LRUcache) Get(key string) (value interface{}, ok bool) {
	var el *list.Element

	c.lock.RLock()
	defer c.lock.RUnlock()

	if el, ok = c.items[key]; !ok {
		return
	}

	c.priorityQueue.MoveToFront(el)
	value = el.Value.(*Item).Value

	return
}

func (c *LRUcache) Remove(key string) (ok bool) {
	var el *list.Element

	c.lock.Lock()
	defer c.lock.Unlock()

	if el, ok = c.items[key]; !ok {
		return true
	}

	c.priorityQueue.Remove(el)
	delete(c.items, key)

	return
}

func (c *LRUcache) getItems() map[string]*list.Element {
	return c.items
}

func (c *LRUcache) getPriorityQueue() *list.List {
	return c.priorityQueue
}
