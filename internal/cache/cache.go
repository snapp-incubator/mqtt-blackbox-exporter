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
