package main

import (
	"database/sql"
	"log"

	"github.com/haunt98/invest-go/internal/cli"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./sql/data.sqlite3")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}

	app, err := cli.NewApp(db)
	if err != nil {
		log.Fatalln(err)
	}

	app.Run()
}
