package module13__06_03

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestLRUcache_Add(t *testing.T) {
	itemsToCache := []int{0, 1, 2}

	tests := map[string]struct {
		expectedItems []int
		expectedQueue []int
		capacity      int
	}{
		"add_less_than_capacity": {
			expectedItems: []int{0, 1, 2},
			expectedQueue: []int{2, 1, 0},
			capacity:      4,
		},
		"add_equal_to_capacity": {
			expectedItems: []int{0, 1, 2},
			expectedQueue: []int{2, 1, 0},
			capacity:      3,
		},
		"add_more_than_capacity": {
			expectedItems: []int{1, 2},
			expectedQueue: []int{2, 1},
			capacity:      2,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := NewLRUCache(tc.capacity)

			for _, v := range itemsToCache {
				cache.Add(strconv.Itoa(v), v)
			}

			items := cache.getItems()
			assertItems(t, tc.expectedItems, items)

			queue := cache.getPriorityQueue()
			assertQueue(t, tc.expectedQueue, queue)
		})
	}
}

func TestLRUcache_GetPrioritize(t *testing.T) {
	capacity := 2
	cache := NewLRUCache(capacity)

	cache.Add("0", 0)
	cache.Add("1", 1)

	v, _ := cache.Get("0")
	assert.Equalf(t, v, 0, "value is not equal to 0")
	cache.Add("2", 2)

	assertItems(t, []int{2, 0}, cache.getItems())
	assertQueue(t, []int{2, 0}, cache.getPriorityQueue())
}

func TestLRUcache_Get(t *testing.T) {
	itemsToCache := []int{0, 1}
	capacity := 2

	tests := map[string]struct {
		getKey         string
		expectedVal    interface{}
		expectedResult bool
	}{
		"get_existing_key": {
			getKey:         "1",
			expectedVal:    1,
			expectedResult: true,
		},
		"get_not_existing_key": {
			getKey:         "2",
			expectedVal:    nil,
			expectedResult: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := NewLRUCache(capacity)

			for _, v := range itemsToCache {
				cache.Add(strconv.Itoa(v), v)
			}

			v, r := cache.Get(tc.getKey)
			assert.Equalf(t, tc.expectedVal, v, "expected value %#v not equal to actual value %#v", tc.expectedVal, v)
			assert.Equalf(t, tc.expectedResult, r, "expected result %#v not equal to actual result %#v", tc.expectedResult, r)
		})
	}
}

func TestLRUcache_Remove(t *testing.T) {
	capacity := 3

	tests := map[string]struct {
		itemsToCache  []int
		removeKey     string
		expectedItems []int
		expectedQueue []int
	}{
		"remove_last_added": {
			itemsToCache:  []int{0, 1, 2},
			removeKey:     "2",
			expectedItems: []int{0, 1},
			expectedQueue: []int{1, 0},
		},
		"remove_first_added": {
			itemsToCache:  []int{0, 1, 2},
			removeKey:     "0",
			expectedItems: []int{1, 2},
			expectedQueue: []int{2, 1},
		},
		"remove_middle_added": {
			itemsToCache:  []int{0, 1, 2},
			removeKey:     "1",
			expectedItems: []int{0, 2},
			expectedQueue: []int{2, 0},
		},
		"empty_result": {
			itemsToCache:  []int{0},
			removeKey:     "0",
			expectedItems: []int{},
			expectedQueue: []int{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			cache := NewLRUCache(capacity)

			for _, v := range tc.itemsToCache {
				cache.Add(strconv.Itoa(v), v)
			}

			cache.Remove(tc.removeKey)

			items := cache.getItems()
			assertItems(t, tc.expectedItems, items)

			queue := cache.getPriorityQueue()
			assertQueue(t, tc.expectedQueue, queue)
		})
	}
}

func assertItems(t *testing.T, expectedItems []int, actualItems map[string]*list.Element) {
	assert.Len(t, actualItems, len(expectedItems))

	for _, ev := range expectedItems {
		k := strconv.Itoa(ev)
		v, ok := actualItems[k]
		assert.Truef(t, ok, "expected key %s not found\n", k)
		av := v.Value.(*Item).Value
		assert.Equalf(t, ev, av, "expected %#v not equal to actual %#v\n", ev, av)
	}
}

func assertQueue(t *testing.T, expectedQueue []int, actualQueue *list.List) {
	assert.Equalf(t, actualQueue.Len(), len(expectedQueue), "queue length %i not equal to cache items length %i\n", actualQueue.Len(), len(expectedQueue))

	eQ := actualQueue.Front()
	for _, v := range expectedQueue {
		kQ := eQ.Value.(*Item).Key
		vQ := eQ.Value.(*Item).Value
		k := strconv.Itoa(v)
		assert.Equalf(t, kQ, k, "queue element key %s not equal to corresponding cache element value %s\n", kQ, k)
		assert.Equalf(t, vQ, v, "queue element value %#v not equal to corresponding cache element value %#v\n", vQ, v)

		eQ = eQ.Next()
	}
}
