run:
	go run main.go --file ./zsh_history

run-debug:
	DEBUG=true $(MAKE) run

test:
	go test ./...

bench:
	go test ./... -benchmem -bench ^Bench

build:
	go build -ldflags "-w -s"

install: build
	mv zsh-history-clear ${SYNC_PATH}/bin
