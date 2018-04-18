package main

import (
	"bufio"
	"os"

	logrus "github.com/sirupsen/logrus"
)

var file = os.Getenv("HOME") + "/.zsh_history"

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

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Fatal(err)
	}

	var myresulter Resulter
	myresulter.ProcessSlice(lines)
	myresulter.WriteFile()
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	do()
}
