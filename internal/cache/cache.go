package cache

import (
	"sync"
	"time"
)

// Item stores information about a pinging.
type Item struct {
	Start  time.Time
	Status bool
}

// Cache is our application cache for clients.
type Cache struct {
	list sync.Map
}

func (c *Cache) Init() {
	c.list = sync.Map{}
}

func (c *Cache) Push(id int, start time.Time) {
	c.list.Store(id, Item{
		Start:  start,
		Status: true,
	})
}

func (c *Cache) Pull(id int) Item {
	value, _ := c.list.Load(id)
	item, _ := value.(Item)

	return item
}
