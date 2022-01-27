package main

import (
	"database/sql"
	"log"
	"time"

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

	// Shout out to Sai Gon, Viet Nam
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatalln(err)
	}

	app, err := cli.NewApp(db, location)
	if err != nil {
		log.Fatalln(err)
	}

	app.Run()
}
