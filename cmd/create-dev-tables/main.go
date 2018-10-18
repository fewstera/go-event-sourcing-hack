// The create-dev-tables command is ran only on developers machines, it is used to create
// the tables on the dockerised SQL database.
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const createSQLFile = "/app/create.sql"

func main() {
	fmt.Println("Creating database tables if the don't already exist")
	db := initDb()
	query := readCreateSQLFile()
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println("Finished creating tables")
}

func readCreateSQLFile() string {
	createSQL, err := ioutil.ReadFile(createSQLFile)
	if err != nil {
		panic(err)
	}
	return string(createSQL)
}

func initDb() *sql.DB {
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/events")
	if err != nil {
		panic(err)
	}

	fmt.Println("Waiting for db to start up...")
	dbConnected := false
	for i := 0; i < 20; i++ {
		err = db.Ping()
		fmt.Printf(".")
		if err == nil {
			dbConnected = true
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("")

	if !dbConnected {
		panic(err)
	}

	return db
}
