run:
	reset && go run main.go counter.go resulter.go

run-debug:
	reset && DEBUG=true go run main.go counter.go resulter.go

test:
	go test ./...

bench:
	go test ./... -benchmem -bench ^Bench
