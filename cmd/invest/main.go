package main

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/haunt98/invest-go/internal/cli"
	"github.com/make-go-great/xdg-go"
	_ "github.com/mattn/go-sqlite3"
)

const dataFilename = "data.sqlite3"

func main() {
	if err := os.MkdirAll(getDataDirPath(), 0o755); err != nil {
		log.Fatalln(err)
	}

	shouldInitDatabase := false
	if _, err := os.Stat(getDataFilePath()); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			shouldInitDatabase = true
		} else {
			log.Fatalln(err)
		}
	}

	db, err := sql.Open("sqlite3", getDataFilePath())
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

// Should be ~/.local/share/invest
func getDataDirPath() string {
	return filepath.Join(xdg.GetDataHome(), cli.Name)
}

// Should be ~/.local/share/invest/data.sqlite3
func getDataFilePath() string {
	return filepath.Join(xdg.GetDataHome(), cli.Name, dataFilename)
}
