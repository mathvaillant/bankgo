build: 
	@go build -o bin/bankgo

run: build
	@./bin/bankgo

test:
	@go test -v ./...
