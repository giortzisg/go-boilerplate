build:
	cd ./deploy/build && docker compose up -d && cd ../../

run:
	go run cmd/server/main.go -config local

test:
	go test -v ./...

test-cover:
	go test -v -cover ./...

lint:
	golangci-lint run ./...
