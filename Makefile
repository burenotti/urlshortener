include .env

EXECUTABLE=app
WINDOWS=$(EXECUTABLE)_windows_amd64.exe
LINUX=$(EXECUTABLE)_linux_amd64

.PHONY:
all: build test clean

windows:
	go build -o ./bin/$(WINDOWS) './cmd/main.go'

unix:
	go build -o ./bin/$(LINUX) './cmd/main.go'

build: windows unix

dev:
ifeq ($(OS), Windows_NT)
	air -c .air.win.toml
else
	air -c .air.toml
endif

run:
	go run ./cmd/main.go

test:
	go test -v -race -count=1 ./...

cover:
	go test -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html="coverage.out"
	rm coverage.out

clean:
	rm ./bin/$(WINDOWS) ./bin/$(LINUX)