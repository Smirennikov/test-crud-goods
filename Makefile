run:
	go run ./cmd/crud $(mode)

migratePostgres:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest &&\
	migrate -database postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DBNAME)?sslmode=disable -path migrations/postgres -verbose $(mode)

migrateClickhouse:
	go install -tags 'clickhouse' github.com/golang-migrate/migrate/v4/cmd/migrate@latest &&\
	migrate -database clickhouse://$(CLICKHOUSE_USER):$(CLICKHOUSE_PASSWORD)@$(CLICKHOUSE_HOST):$(CLICKHOUSE_PORT)/$(CLICKHOUSE_DB) -path migrations/clickhouse -verbose $(mode)

migrate:
	make migratePostgres && make migrateClickhouse