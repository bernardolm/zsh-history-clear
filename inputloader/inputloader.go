package inputloader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/bernardolm/zsh-history-clear/executor"
	log "github.com/sirupsen/logrus"
)

type ex struct {
	wg *sync.WaitGroup
}

func (e *ex) Get() []byte {
	defer e.wg.Done()
	e.wg.Add(1)

	in := os.Stdin

	stdinStat, err := in.Stat()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		panic(err)
	}

	if (stdinStat.Mode() & os.ModeCharDevice) == 0 {
		scanner := bufio.NewScanner(in)
		for {
			log.Debugf(".")
			if !scanner.Scan() {
				break
			}

			b := scanner.Bytes()
			log.Debugf("inputloader.Get.Bytes(): %v\n", string(b))

			if len(b) == 0 {
				continue
			}

			return b
		}
	} else {
		log.Info("Enter file path: ")
		var filePath string
		fmt.Scanf("%s", &filePath)
		log.Info(filePath)
	}

	return nil
}

func (e *ex) Put(_ []byte) {
	defer e.wg.Done()
	e.wg.Add(1)
}

func New(wg *sync.WaitGroup) executor.Executor {
	e := ex{
		wg: wg,
	}
	return &e
}
