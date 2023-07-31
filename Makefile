build:
	@go build -o bin/clock main.go

build-win:
	@GOOS=windows GOARCH=amd64 go build -o bin/clock.exe main.go

run: build
	@./bin/clock

