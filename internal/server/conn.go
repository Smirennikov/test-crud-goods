package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/fatih/color"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/nats-io/nats.go"
)

func redis_connect(db int, host, port, password string) (*redis.Client, error) {
	addr := net.JoinHostPort(host, port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	log.Println(color.GreenString("Redis successfully connected %s", addr))
	return rdb, nil
}

func pg_connect(host, port, user, password, dbname string) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	conn, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	log.Println(color.GreenString("Postgres pool successfully connected %s:%s", host, port))
	return conn, nil
}

func clickhouse_connect(host, port, database, username, password string) (driver.Conn, error) {
	var (
		ctx       = context.Background()
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{host + ":" + port},
			// Auth: clickhouse.Auth{
			// 	Database: database,
			// 	Username: username,
			// 	Password: password,
			// },
			Debugf: func(format string, v ...interface{}) {
				fmt.Printf(format, v)
			},
		})
	)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		}
		return nil, err
	}

	log.Println(color.GreenString("Clickhouse successfully connected %s:%s", host, port))
	return conn, nil
}

func nats_connect() (*nats.Conn, error) {
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}
	log.Println(color.GreenString("Nats successfully connected %s", nats.DefaultURL))
	return conn, nil
}
