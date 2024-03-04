# Test task
The app implements basic operations with data known as CRUD. Here uses Postgres as the main database, Redis for cache, NATS and ClickHouse for processing logs.

## How to run dev environment for app
`docker-compose -f docker-compose.dev.yml up -d`

## How to run migrations
1. `export $(grep -v '^#' dev.env | xargs)`
2. `make migrate mode=up/down`

## How to run app in dev mode
`make run mode=dev`

## How to run lint
`golangci-lint run`
