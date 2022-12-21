fmt:
	go fmt cmd/gedis/main.go

run:
	go run cmd/gedis/main.go

build:
	go build -o gedis cmd/gedis/main.go

test:
	curl -X POST http://localhost:8000/set -H 'Content-Type: application/json' -d '{"user":"chilts","password":"password"}'
