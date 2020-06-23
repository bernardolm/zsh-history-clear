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
	var l []string
	for sc.Scan() {
		l = append(l, sc.Text())
	}
	if err := sc.Err(); err != nil {
		logrus.WithError(err).Panic(err)
	}
	return l
}

func splitLines(l []string) map[string]string {
	m := map[string]string{}
	rs := `^(:\s\d+:\d+;)(.*)$`
	re := regexp.MustCompile(rs)
	if re == nil {
		logrus.Panicf("regex %s not compile", rs)
	}
	sort.Strings(l)
	for k := range l {
		g := re.FindAllStringSubmatch(l[k], -1)
		if len(g) == 0 {
			logrus.WithField("line", l[k]).Panic("line can't match regex")
		}
		cmd := strings.TrimSpace(g[0][2])
		m[cmd] = fmt.Sprintf("%s%s", g[0][1], cmd)
	}
	logrus.Infof("turn %d lines into %d", len(l), len(m))
	return m
}

func writeFile(m map[string]string) {
	var b bytes.Buffer
	for k := range m {
		b.WriteString(m[k])
		b.WriteString("\n")
	}
	if err := ioutil.WriteFile(outputFilePath, b.Bytes(), 0664); err != nil {
		logrus.WithError(err).Panic(err)
	}
}

func main() {
	initLogger()
	writeFile(splitLines(parseLines(readFile())))
}
