package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResulter(t *testing.T) {
	var resulter = new(Resulter)

	data := []string{
		"Smells Like Teen Spirit",
		"In Bloom",
		"Come as You Are",
		"Breed",
		"Breed",
		"Lithium",
		"Polly",
		"Territorial Pissings",
		"Drain You",
		"Lounge Act",
		"Lounge Act",
		"Stay Away",
		"On a Plain",
		"Something in the Way",
		"Endless, Nameless",
		"Endless, Nameless",
		"Endless, Nameless",
		"Endless, Nameless",
	}

	splitter := func(s string) (string, string, bool) {
		return s, s, true
	}

	resulter.ProcessSlice(data, splitter)

	assert.NotEmpty(t, resulter.Result())
	assert.Len(t, resulter.Result(), 13)

	assert.Equal(t, data[0], resulter.Result()[data[0]])
	assert.Equal(t, data[1], resulter.Result()[data[1]])
	assert.Equal(t, data[5], resulter.Result()[data[5]])
	assert.Equal(t, data[10], resulter.Result()[data[10]])
	assert.Equal(t, data[12], resulter.Result()[data[12]])
}
