package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"sort"

	logrus "github.com/sirupsen/logrus"
)

// TODO: With less limit, don't match repetead
const limit int = 250000

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
	s      []string
}

func (r *resulter) Len() int {
	return len(r.result)
}

func (r *resulter) Less(i, j int) bool {
	return r.s[i] < r.s[j]
}

func (r *resulter) Swap(i, j int) {
	r.s[i], r.s[j] = r.s[j], r.s[i]
}

func sortedKeys(m map[string]string) []string {
	sm := new(resulter)
	sm.result = m
	sm.s = make([]string, len(m))
	i := 0
	for key := range m {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
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
	for _, v := range sortedKeys(r.result) {
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
