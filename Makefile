build: 
	@go build -o ./bin/main cmd/*.go

test:
	@go test -v ./...

run: 
	@air 