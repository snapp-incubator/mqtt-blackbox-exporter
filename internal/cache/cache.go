package cache

import (
	"sync"
	"time"
)

// Item stores information about a pinging.
type Item struct {
	Start time.Time
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
		Start: start,
	})
}

func (c *Cache) Pull(id int) (Item, bool) {
	value, has := c.list.LoadAndDelete(id)
	item, _ := value.(Item)

	return item, has
}
