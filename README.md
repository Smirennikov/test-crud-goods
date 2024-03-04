# test task for create CRUD with writing logs to ClickHouse through NATS

## How to run dev environment for app
`docker-compose -f docker-compose.dev.yml up -d`

## How to run migrations
1. `export $(grep -v '^#' dev.env | xargs)`
2. `make migrate mode=up/down`

## How to run app in dev mode
`make run mode=dev`

## How to run lint
`golangci-lint run`
