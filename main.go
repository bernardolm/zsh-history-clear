package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"

	logrus "github.com/sirupsen/logrus"
)

// TODO: With less limit, don't match repetead
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
			logrus.WithField("entry", v).WithField("value", value).Debug("ignoring repeated")
		}
	}

	mycounter.reset()
	lines = &[]string{}
}

func (r resulter) writeFile() {
	var buffer bytes.Buffer
	for _, v := range r.result {
		logrus.WithField("value", v).Debug("writing line to file")
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}

	err := ioutil.WriteFile(file+"_new", buffer.Bytes(), 0664)
	if err != nil {
		logrus.WithError(err).Fatal(err)
	}
}

var myresulter resulter

func do() {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		logrus.WithError(err).Panic("filepath not exist")
	}

	file, err := os.Open(file)
	if err != nil {
		logrus.Fatal(err)
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
		logrus.WithError(err).Fatal(err)
	}

	myresulter.writeFile()
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	do()
}
