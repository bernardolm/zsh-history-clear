package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResulter(t *testing.T) {
	var r = new(Resulter)

	lines := []string{
		"Smells Like Teen Spirit",
		"In Bloom",
		"Come as You Are",
		"Breed",
		"Lithium",
		"Polly",
		"Territorial Pissings",
		"Drain You",
		"Lounge Act",
		"Stay Away",
		"On a Plain",
		"Something in the Way",
		"Endless, Nameless",
	}

	assert.NotNil(t, r)
	assert.NotEmpty(t, lines)
}
