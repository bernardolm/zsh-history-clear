package ohmpatel

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

func (e *ex) do(file *os.File) error {
	filestat, err := file.Stat()
	if err != nil {
		return errors.New("Could not able to get the file stat")
	}

	fileSize := filestat.Size()
	offset := fileSize - 1
	lastLineSize := 0

	for {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, offset)
		if err != nil {
			log.Errorf("Error reading file ", err)
			break
		}
		char := string(b[0])
		if char == "\n" {
			break
		}
		offset--
		lastLineSize += n
	}

	lastLine := make([]byte, lastLineSize)
	_, err = file.ReadAt(lastLine, offset+1)

	if err != nil {
		return fmt.Errorf("Could not able to read last line with offset %d  and lastline size %d", offset, lastLineSize)
	}

	return e.process(file)
}

func (e *ex) process(f *os.File) error {
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)

	// var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				panic(err)
			}
			return err
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		// wg.Add(1)
		e.wg.Add(1)
		go func() {
			e.processChunk(buf, &linesPool, &stringPool)
			// wg.Done()
			e.wg.Done()
		}()
	}

	// wg.Wait()
	return nil
}

func (e *ex) processChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool) {
	// var wg2 sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "\n")

	stringPool.Put(logs)

	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {
		// wg2.Add(1)
		e.wg.Add(1)

		func(start int, end int) {
			// defer wg2.Done() //to avaoid deadlocks
			defer e.wg.Done()

			for i := start; i < end; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}

				go e.next.Put([]byte(text))
			}

		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}

	// wg2.Wait()
	logsSlice = nil
}
