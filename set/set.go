package set

import (
	"fmt"
)

//MakeSet initialize the set
func makeSet() *customSet {
	return &customSet{
		Container: make(map[string]struct{}),
	}
}

type customSet struct {
	Container map[string]struct{}
}

func (c *customSet) Exists(key string) bool {
	_, exists := c.Container[key]
	return exists
}

func (c *customSet) Add(key string) {
	c.Container[key] = struct{}{}
}

func (c *customSet) Remove(key string) error {
	_, exists := c.Container[key]
	if !exists {
		return fmt.Errorf("Remove Error: Item doesn't exist in set")
	}
	delete(c.Container, key)
	return nil
}

func (c *customSet) Size() int {
	return len(c.Container)
}
