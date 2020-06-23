run:
	go run main.go --file ./zsh_history

run-debug:
	DEBUG=true $(MAKE) run

test:
	go test ./...

bench:
	go test ./... -benchmem -bench ^Bench
