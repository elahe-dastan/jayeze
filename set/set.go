package set

import (
	"fmt"
)

//MakeSet initialize the set
func MakeSet() *customSet {
	return &customSet{
		Container: make(map[int]struct{}),
	}
}

type customSet struct {
	Container map[int]struct{}
}

func (c *customSet) Exists(key int) bool {
	_, exists := c.Container[key]
	return exists
}

func (c *customSet) Add(key int) {
	c.Container[key] = struct{}{}
}

func (c *customSet) Remove(key int) error {
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
