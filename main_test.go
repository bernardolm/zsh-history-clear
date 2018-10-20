package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/stretchr/testify/assert"
)

func TestSplitZshHistoryKeyValue(t *testing.T) {
	input := ": 1539986823:0;ls -lah"
	expected := "ls -lah"

	actual1, actual2, actualBool := splitZshHistoryKeyValue(input)
	assert.Equal(t, expected, actual1)
	assert.Equal(t, input, actual2)
	assert.True(t, actualBool)

	input = "x 1539986823:0;ls -lah"

	actual1, actual2, actualBool = splitZshHistoryKeyValue(input)
	assert.Equal(t, "", actual1)
	assert.Equal(t, input, actual2)
	assert.False(t, actualBool)

	input = ": 1539986823:0"

	actual1, actual2, actualBool = splitZshHistoryKeyValue(input)
	assert.Equal(t, "", actual1)
	assert.Equal(t, input, actual2)
	assert.False(t, actualBool)
}

func createFileTest() string {
	data, _ := ioutil.ReadFile("zsh_history")
	dst := "zsh_history_" + time.Now().Format("20060102030405.000")
	ioutil.WriteFile(dst, data, 0644)
	return dst
}

func benchmarkDo(i int, b *testing.B) {
	logrus.SetLevel(logrus.ErrorLevel)
	for n := 0; n < 100; n++ {
		filePath = createFileTest()
		do()
		os.Remove(filePath)
	}
}

func BenchmarkDo1(b *testing.B)           { benchmarkDo(1, b) }
func BenchmarkDo100(b *testing.B)         { benchmarkDo(100, b) }
func BenchmarkDo1000(b *testing.B)        { benchmarkDo(1000, b) }
func BenchmarkDo10000(b *testing.B)       { benchmarkDo(10000, b) }
func BenchmarkDo100000(b *testing.B)      { benchmarkDo(100000, b) }
func BenchmarkDo1000000(b *testing.B)     { benchmarkDo(1000000, b) }
func BenchmarkDo10000000(b *testing.B)    { benchmarkDo(10000000, b) }
func BenchmarkDo100000000(b *testing.B)   { benchmarkDo(100000000, b) }
func BenchmarkDo1000000000(b *testing.B)  { benchmarkDo(1000000000, b) }
func BenchmarkDo10000000000(b *testing.B) { benchmarkDo(10000000000, b) }
