package main

type Counter struct {
	count      int
	totalCount int
}

func (c *Counter) Plus() {
	c.count++
	c.totalCount++
}

func (c *Counter) Reset() {
	c.count = 0
}

func (c Counter) Position() int {
	return c.count
}

func (c Counter) Total() int {
	return c.totalCount
}
