all: run


install:
	rm -rf pebble
	go build ./pebble.go
	sudo cp pebble /usr/bin/pebble

uninstall:
	sudo rm -rf /usr/bin/pebble

build:
	rm -rf pebble
	go build -o pebble -ldflags "-s -w" -trimpath pebble.go

run:
	go run ./pebble.go schema migrate
