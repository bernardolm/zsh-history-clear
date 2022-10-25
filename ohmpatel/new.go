package ohmpatel

import (
	"errors"
	"io"
	"os"
	"sync"

	"github.com/bernardolm/zsh-history-clear/executor"
	log "github.com/sirupsen/logrus"
)

type ex struct {
	input executor.Executor
	next  executor.Executor
	wg    *sync.WaitGroup
}

func (e *ex) Get() []byte {
	defer e.wg.Done()
	e.wg.Add(1)

	b := e.input.Get()

	f, err := os.Open(string(b))
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}

		panic(err)
	}

	if err := e.do(f); err != nil {
		log.Panic(err)
	}

	return nil
}

func (e *ex) Put(_ []byte) {
	defer e.wg.Done()
	e.wg.Add(1)
}

func New(input executor.Executor, next executor.Executor, wg *sync.WaitGroup) executor.Executor {
	e := ex{
		input: input,
		next:  next,
		wg:    wg,
	}
	return &e
}
