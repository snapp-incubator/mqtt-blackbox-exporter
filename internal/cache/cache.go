package cache

import (
	"time"
)

// Item stores information about a pinging.
type Item struct {
	Start  time.Time
	Status bool
}

// Cache is our application cache for clients.
type Cache struct {
	list map[int]Item
}

func (c *Cache) Init() {
	c.list = make(map[int]Item)
}

func (c *Cache) Push(id int, start time.Time) {
	c.list[id] = Item{
		Start:  start,
		Status: false,
	}
}

func (c *Cache) Pull(id int) Item {
	return c.list[id]
}
