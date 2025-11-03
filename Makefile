.PHONY: build run

APP_NAME = build/dbash

build:
	mkdir -p build
	go clean -cache -modcache
	go mod tidy
	go build -o $(APP_NAME) main.go dbash.go executor.go

run: build
	clear
	./$(APP_NAME) --local=teleport.dbash

run-debug: build
	clear
	./$(APP_NAME) --local=teleport.dbash --debug
