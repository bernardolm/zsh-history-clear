package main

import (
	"bufio"
	"flag"
	"os"

	logrus "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var filePath string

func splitZshHistoryKeyValue(s string) (string, string, bool) {
	if len(s) <= 15 || s[0:1] != ":" {
		return "", s, false
	}

	return s[15:len(s)], s, true
}

func do() {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logrus.WithError(err).Panicf("filepath %s not exist", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		logrus.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Fatal(err)
	}

	var myresulter Resulter
	myresulter.ProcessSlice(lines, splitZshHistoryKeyValue)
	myresulter.WriteFile()
}

func init() {
	logrus.SetOutput(os.Stdout)

	if viper.GetBool("DEBUG") {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func main() {
	file := flag.String("file", os.Getenv("HOME")+"/.zsh_history", "file path")
	filePath = *file

	do()
}
