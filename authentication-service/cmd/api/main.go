package main

import (
	"authentication/cmd/api/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const (
	webPort = "80"
)

var tryCounts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting authentication service")

	// connect to databse
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to postgres")
	}

	// setup config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	src := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := src.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready...")
			tryCounts++
		} else {
			log.Println("connected to postgres!")
			return connection
		}

		if tryCounts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("try to connect again in 2s...")
		time.Sleep(2 * time.Second)
		continue
	}
}
