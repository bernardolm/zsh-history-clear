package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	var c = new(Counter)

	assert.Equal(t, 0, c.Position())
	assert.Equal(t, 0, c.Total())

	for i := 0; i < 12; i++ {
		c.Plus()
	}

	assert.Equal(t, 12, c.Position())
	assert.Equal(t, 12, c.Total())

	for i := 0; i < 7; i++ {
		c.Plus()
	}

	assert.Equal(t, 19, c.Position())
	assert.Equal(t, 19, c.Total())

	c.Reset()

	assert.Equal(t, 0, c.Position())
	assert.Equal(t, 19, c.Total())
}
