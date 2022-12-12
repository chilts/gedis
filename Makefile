fmt:
	go fmt cmd/gedis/main.go

run:
	go run cmd/gedis/main.go

build:
	go build -o gedis cmd/gedis/main.go
