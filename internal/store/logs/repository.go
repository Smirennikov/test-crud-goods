package logs

import "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

type logs struct {
	db driver.Conn
}

func New(db driver.Conn) *logs {
	return &logs{
		db: db,
	}
}
