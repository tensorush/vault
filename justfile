run:
    go mod verify
    go run ./cmd/vault-bot/main.go

test:
	go test -race -vet=off ./internal/database/sqlite ./internal/vault

docker:
    docker-compose up --build

docker-test:
	go test -race -vet=off ./...
