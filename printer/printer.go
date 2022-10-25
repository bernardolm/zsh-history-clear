package printer

import (
	"sync"

	"github.com/bernardolm/zsh-history-clear/executor"
	log "github.com/sirupsen/logrus"
)

type ex struct {
	wg *sync.WaitGroup
	i  int
}

func (e *ex) Get() []byte {
	defer e.wg.Done()
	e.wg.Add(1)

	return nil
}

func (e *ex) Put(b []byte) {
	defer e.wg.Done()
	e.wg.Add(1)

	if len(b) > 0 {
		e.i += 1
	}

	log.Debugf("%d", e.i)
}

func New(wg *sync.WaitGroup) executor.Executor {
	e := ex{
		wg: wg,
	}
	return &e
}
