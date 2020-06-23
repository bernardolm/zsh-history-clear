package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/cheggaaa/pb/v3"
	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var outputFilePath string

func initLogger() {
	logrus.SetOutput(os.Stdout)
	if viper.GetBool("DEBUG") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func readFile() *os.File {
	fp := flag.String("file", "./zsh_history", "file path")
	flag.Parse()
	if fp == nil {
		logrus.Panic("filepath not exist")
	}
	if _, err := os.Stat(*fp); os.IsNotExist(err) {
		logrus.WithError(err).Panicf("filepath %s not exist", *fp)
	}
	outputFilePath = *fp
	f, err := os.Open(*fp)
	if err != nil {
		logrus.WithError(err).Panic(err)
	}
	return f
}

func parseLines(f *os.File) []string {
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	l := []string{}
	for sc.Scan() {
		if len(sc.Text()) > 0 {
			l = append(l, sc.Text())
		}
	}
	if err := sc.Err(); err != nil {
		logrus.WithError(err).Panic(err)
	}
	return l
}

func uniqueLines(l []string) []string {
	m := make(map[string]int, len(l))
	rs := `^(:\s\d+:\d+;)(.*)$`
	re := regexp.MustCompile(rs)
	if re == nil {
		logrus.Panicf("regex %s not compile", rs)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(l)))
	lo := make([]string, 0, len(l))
	fmt.Println("processing...")
	bar := pb.StartNew(len(l))
	for k := range l {
		g := re.FindAllStringSubmatch(l[k], -1)
		if len(g) == 0 {
			logrus.WithField("line", l[k]).Panic("line can't match regex")
		}
		cmd := strings.TrimSpace(g[0][2])
		if _, ok := m[cmd]; !ok {
			lo = append(lo, fmt.Sprintf("%s%s", g[0][1], cmd))
			m[cmd] = 1
		}
		bar.Increment()
	}
	bar.Finish()
	logrus.Infof("turn %d lines into %d", len(l), len(lo))
	sort.Strings(lo)
	return lo
}

func writeFile(l []string) {
	var b bytes.Buffer
	fmt.Println("writing...")
	bar := pb.StartNew(len(l))
	for k := range l {
		b.WriteString(l[k])
		b.WriteString("\n")
		bar.Increment()
	}
	bar.Finish()
	if err := ioutil.WriteFile(outputFilePath, b.Bytes(), 0664); err != nil {
		logrus.WithError(err).Panic(err)
	}
}

func main() {
	initLogger()
	writeFile(
		uniqueLines(
			parseLines(
				readFile())))
}
