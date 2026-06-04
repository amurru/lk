.PHONY: build test install lint clean

build:
	go build -o lk ./cmd/lk/

test:
	go test ./...

install: build
	cp lk /usr/local/bin/

lint:
	go vet ./...

clean:
	rm -f lk
