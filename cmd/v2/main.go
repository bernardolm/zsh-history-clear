package main

import (
	"sync"
	"time"

	"github.com/bernardolm/zsh-history-clear/inputloader"
	"github.com/bernardolm/zsh-history-clear/ohmpatel"
	"github.com/bernardolm/zsh-history-clear/printer"
	log "github.com/sirupsen/logrus"
)

var (
	wg sync.WaitGroup
)

func main() {
	log.SetLevel(log.DebugLevel)

	start := time.Now()

	input := inputloader.New(&wg)
	next := printer.New(&wg)
	op := ohmpatel.New(input, next, &wg)

	op.Get()

	wg.Wait()

	log.Infof("Time taken: %s", time.Since(start))
}
