package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/haunt98/invest-go/internal/cli"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	shouldInitDatabase := false
	if _, err := os.Stat("./sql/data.sqlite3"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			shouldInitDatabase = true
		} else {
			log.Fatalln(err)
		}
	}

	db, err := sql.Open("sqlite3", "./sql/data.sqlite3")
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	// Shout out to Sai Gon, Viet Nam
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatalln(err)
	}

	app, err := cli.NewApp(db, shouldInitDatabase, location)
	if err != nil {
		log.Fatalln(err)
	}

	app.Run()
}
