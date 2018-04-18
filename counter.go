package main

type Counter struct {
	count int
	limit int
	total int
}

func (c *Counter) Plus() {
	c.count++
	c.total++
}

func (c *Counter) Reset() {
	c.count = 0
}

func (c Counter) Position() int {
	return c.count
}

func (c Counter) Total() int {
	return c.total
}

func (c Counter) NotReached() bool {
	return c.count < c.limit
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
