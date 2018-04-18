package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const limit int = 250

var file = os.Getenv("HOME") + "/.zsh_history"

type counter struct {
	count      int
	totalCount int
}

func (c *counter) plus() {
	c.count++
	c.totalCount++
}

func (c *counter) reset() {
	c.count = 0
}

func (c *counter) position() int {
	return c.count
}

func (c *counter) total() int {
	return c.totalCount
}

type resulter struct {
	sync.Mutex

	result map[string]string
}

func (r *resulter) addData(lines *[]string, mycounter *counter) {
	r.Lock()
	defer r.Unlock()

	if r.result == nil {
		r.result = make(map[string]string)
	}
	for _, v := range *lines {
		if len(v) < 16 {
			continue
		}
		f := v[0:1]
		value := v[16:len(v)]
		if _, ok := r.result[value]; !ok && f == ":" {
			r.result[value] = v
		} else {
			fmt.Printf("ignoring entry repeated %s\n", v)
		}
	}

	mycounter.reset()
	lines = &[]string{}
}

func (r resulter) writeFile() {
	var buffer bytes.Buffer
	for _, v := range r.result {
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}

	err := ioutil.WriteFile(file, buffer.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

var myresulter resulter

func do() {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic("filepath not exist")
	}

	file, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	mycounter := counter{}
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		mycounter.plus()

		if mycounter.position() == limit {
			go myresulter.addData(&lines, &mycounter)
		}
	}

	myresulter.addData(&lines, &mycounter)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	myresulter.writeFile()
}

func main() {
	do()
}
