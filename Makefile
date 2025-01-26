build:
	cd ./deploy/build && CONFIG_ENV=dev docker compose up -d --build && cd ../../

run:
	go run cmd/migration/main.go -config config/local.yaml && go run cmd/server/main.go -config config/local.yaml

test:
	go test -v ./...

test-cover:
	go test -v -cover ./...

lint:
	golangci-lint run ./...
