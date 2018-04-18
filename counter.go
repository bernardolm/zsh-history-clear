package main

type Counter struct {
	count      int
	totalCount int
}

func (c *Counter) plus() {
	c.count++
	c.totalCount++
}

func (c *Counter) reset() {
	c.count = 0
}

func (c Counter) position() int {
	return c.count
}

func (c Counter) total() int {
	return c.totalCount
}
