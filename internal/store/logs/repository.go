package logs

import "github.com/ClickHouse/clickhouse-go/v2/lib/driver"

type Logs struct {
	db driver.Conn
}

func New(db driver.Conn) *Logs {
	return &Logs{
		db: db,
	}
}
