package main

import (
	"bufio"
	"os"

	logrus "github.com/sirupsen/logrus"
)

// TODO: don't matching repetead With limit less than total lines
const limit int = 50000000

var file = os.Getenv("HOME") + "/.zsh_history"

var myresulter Resulter

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

	mycounter := Counter{}
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		mycounter.Plus()

		if mycounter.Position() == limit {
			go myresulter.AddData(&lines, &mycounter)
		}
	}

	myresulter.AddData(&lines, &mycounter)

	if err := scanner.Err(); err != nil {
		logrus.WithError(err).Fatal(err)
	}

	myresulter.WriteFile()
}

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
}

func main() {
	do()
}
