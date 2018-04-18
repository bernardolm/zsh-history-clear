package main

type Counter struct {
	count      int
	limit      int
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

func NewCounter(args ...interface{}) *Counter {
	if len(args) > 0 {
		if limit, ok := args[0].(*int); ok && limit != nil && *limit > 0 {
			return &Counter{
				limit: *limit,
			}
		}
	}

	return &Counter{}
}
