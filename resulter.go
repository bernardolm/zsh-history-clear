package main

import (
	"bytes"
	"io/ioutil"
	"sort"
	"sync"

	logrus "github.com/sirupsen/logrus"
	"vbom.ml/util/sortorder"
)

type KeyValueChecker (func(string) (string, string, bool))

type Resulter struct {
	sync.Mutex

	result map[string]string
	s      []string
}

func (r *Resulter) ProcessSlice(data []string, kvc KeyValueChecker) {
	if data == nil {
		logrus.Error("empty data")
		return
	}

	// TODO: don't matching repetead With limit less than total lines
	limit := 50000000
	if limit > len(data) {
		limit = len(data)
	}

	for position := 0; position <= len(data); position += limit {
		go r.addData(data[position:limit], kvc)
	}
}

func (r *Resulter) addData(data []string, kvc KeyValueChecker) {
	r.Lock()
	defer r.Unlock()

	if r.result == nil {
		r.result = make(map[string]string)
	}

	for _, item := range data {
		key, value, ok := kvc(item)
		if !ok {
			continue
		}

		if _, ok = r.result[key]; !ok {
			r.result[key] = value
		} else {
			logrus.WithField("entry", item).
				WithField("key", key).
				WithField("value", value).
				Debug("ignoring repeated")
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
