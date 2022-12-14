package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"authentication/data"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const PORT = "8081"

var counts int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting authentication service")

	conn := connectToDB()
	if conn == nil {
		log.Fatal("can't connect to postgres...")
	}

	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
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
			counts += 1
		} else {
			log.Println("connected to postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("backing off for 2 seconds")
		time.Sleep(time.Second * 2)
		continue
	}
}
