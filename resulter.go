package main

import (
	"bytes"
	"io/ioutil"
	"sort"
	"sync"

	logrus "github.com/sirupsen/logrus"
	"vbom.ml/util/sortorder"
)

type Resulter struct {
	sync.Mutex

	counter *Counter
	result  map[string]string
	s       []string
}

func (r *Resulter) ProcessSlice(data []string) {
	if data == nil {
		logrus.Error("empty data")
		return
	}

	// TODO: don't matching repetead With limit less than total lines
	r.counter = NewCounter(50000000)

	for _ = range data {
		r.counter.Plus()

		if r.counter.NotReached() {
			go r.addData(data)
		}
	}
}

func (r *Resulter) addData(lines []string) {
	r.Lock()
	defer r.Unlock()

	if r.result == nil {
		r.result = make(map[string]string)
	}

	for _, v := range lines {
		if len(v) <= 15 {
			continue
		}
		first := v[0:1]
		value := v[15:len(v)]
		if _, ok := r.result[value]; !ok && first == ":" {
			r.result[value] = v
		} else {
			logrus.WithField("entry", v).WithField("value", value).Debug("ignoring repeated")
		}
	}
}

func (r *Resulter) Len() int {
	return len(r.result)
}

func (r *Resulter) Less(i, j int) bool {
	return sortorder.NaturalLess(r.s[i], r.s[j])
}

func (r *Resulter) Swap(i, j int) {
	r.s[i], r.s[j] = r.s[j], r.s[i]
}

func (r *Resulter) sortedKeys() []string {
	sm := new(Resulter)
	sm.result = r.result
	sm.s = make([]string, len(r.result))
	i := 0
	for key := range r.result {
		sm.s[i] = key
		i++
	}
	sort.Sort(sm)
	return sm.s
}

func (r *Resulter) WriteFile() {
	var buffer bytes.Buffer
	for _, v := range r.sortedKeys() {
		logrus.WithField("value", v).Debug("writing line to file")
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}

	err := ioutil.WriteFile(file+"_new", buffer.Bytes(), 0664)
	if err != nil {
		logrus.WithError(err).Fatal(err)
	}
}
