all: clean build run

fastrun:
	go run src/main.go --config config.yaml

build:
	go build -o build/app src/main.go

clean: 
	rm -rf ./build
	go clean

run:
	./build/app

