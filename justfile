up *SERVICES:
    docker compose --env-file ./configs/config.env up --build -d {{ SERVICES }}
    docker compose --env-file ./configs/config.env logs

migrate ADDRESS="postgresql://root:secret@0.0.0.0:5432/vault?sslmode=disable" DIRECTION="up":
    migrate -path ./internal/db/migrations/ -database "{{ ADDRESS }}" -verbose {{ DIRECTION }}

migration NAME:
    migrate create -ext sql -dir ./internal/db/migrations/ -seq {{ NAME }}

down:
    docker compose --env-file ./configs/config.env down
