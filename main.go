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

	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var outputFilePath string
var linesRead int

func newProgressBar(size int, title string) *progressbar.ProgressBar {
	return progressbar.NewOptions(size,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetDescription("[cyan]"+title+"[reset]"),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			BarEnd:        "]",
			BarStart:      "[",
			Saucer:        "[green]▉[reset]",
			SaucerHead:    "[green]▁[reset]",
			SaucerPadding: "░",
		}),
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
	)
}

func initLogger() {
	log.SetOutput(os.Stdout)
	if viper.GetBool("DEBUG") {
		log.SetLevel(log.DebugLevel)
	}
}

func readFile() *os.File {
	fp := flag.String("file", "./zsh_history", "file path")
	df := flag.Bool("debug", false, "debug mode")
	flag.Parse()
	if fp == nil {
		log.Panic("filepath not exist")
	}
	if *df {
		log.SetLevel(log.DebugLevel)
	}
	if _, err := os.Stat(*fp); os.IsNotExist(err) {
		log.WithError(err).Panicf("filepath %s not exist", *fp)
	}
	outputFilePath = *fp
	f, err := os.Open(*fp)
	if err != nil {
		log.WithError(err).Panic(err)
	}
	return f
}

func parseLines(f *os.File) []string {
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Split(bufio.ScanLines)
	l := []string{}
	for sc.Scan() {
		if len(sc.Text()) > 0 && sc.Text()[:1] == ":" {
			l = append(l, sc.Text())
		}
	}
	if err := sc.Err(); err != nil {
		log.WithError(err).Panic(err)
	}
	linesRead = len(l)
	return l
}

func uniqueLines(l []string) []string {
	m := make(map[string]int, len(l))
	rs := `^(:\s\d+:\d+;)(.*)$`
	re := regexp.MustCompile(rs)
	if re == nil {
		log.Panicf("regex %s not compile", rs)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(l)))
	lo := make([]string, 0, len(l))
	bar1 := newProgressBar(len(l), "processing...  ")
	for k := range l {
		g := re.FindAllStringSubmatch(l[k], -1)
		if len(g) == 0 {
			log.WithField("line", l[k]).Panic("line can't match regex")
		}
		cmd := strings.TrimSpace(g[0][2])
		if _, ok := m[cmd]; !ok {
			lo = append(lo, fmt.Sprintf("%s%s", g[0][1], cmd))
			m[cmd] = 1
		}
		_ = bar1.Add(1)
	}
	fmt.Println()
	sort.Strings(lo)
	return lo
}

func uniqueLines2(l []string) []string {
	m := make(map[string]int, len(l))
	sort.Sort(sort.Reverse(sort.StringSlice(l)))
	lo := make([]string, 0, len(l))
	bar1 := newProgressBar(len(l), "processing...  ")
	for k := range l {
		if len(l[k]) <= 15 {
			log.WithField("line", l[k]).
				WithField("index", l[k][:15]).
				Error("line without command")
			continue
		}
		log.WithField("line", l[k]).
			WithField("index", l[k][:15]).
			WithField("command", strings.TrimSpace(l[k][15:len(l[k])])).
			Debug("each line")
		cmd := strings.TrimSpace(l[k][15:len(l[k])])
		if _, ok := m[cmd]; !ok {
			lo = append(lo, fmt.Sprintf("%s%s", l[k][:15], cmd))
			m[cmd] = 1
		}
		_ = bar1.Add(1)
	}
	fmt.Println()
	sort.Strings(lo)
	return lo
}

func writeFile(l []string) {
	var b bytes.Buffer
	bar2 := newProgressBar(len(l), "writing file...")
	for k := range l {
		b.WriteString(l[k])
		b.WriteString("\n")
		_ = bar2.Add(1)
	}
	fmt.Println()
	if err := ioutil.WriteFile(outputFilePath, b.Bytes(), 0664); err != nil {
		log.WithError(err).Panic(err)
	}
	fmt.Printf("turn %d lines into %d\n", linesRead, len(l))
}

func main() {
	initLogger()
	writeFile(
		uniqueLines2(
			parseLines(
				readFile())))
}
