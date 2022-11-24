package mapdb

import (
	"fmt"
	"os"
	"strings"

	"github.com/sojebsikder/go-base/src/lib"
	"github.com/sojebsikder/go-base/src/util"
)

// var filename = "db/db2.json"
var dbFileNameDir string = "db"
var dbFileName string

// Interactive cli
func Cli() {
	// first ask for database name
	if dbFileName == "" {
		fmt.Print("Select database: ")
		var dbName string
		fmt.Scanln(&dbName)
		dbFileName = dbFileNameDir + "/" + strings.Trim(dbName, "\n") + ".json"

		_, err := os.Stat(dbFileName)

		if os.IsNotExist(err) {
			// create new database file
			ok := lib.YesNoPrompt("db not exist, create new one? (y/n): ")
			if ok {
				fmt.Print("enter new database name -> ")
				var text string
				fmt.Scan(&text)
				util.CreateDir(dbFileNameDir)
				util.CreateFile(dbFileNameDir + "/" + text + ".json")
				fmt.Println("Database created: " + dbFileName)
			} else {
				fmt.Println("So you want to create later.")
			}

		} else {
			fmt.Println("Database selected: " + dbFileName)
			MapDB()
		}
	}

}

// CLI based simple db operations for store in memory
func MapDB() {

	var cmd string
	var db = map[string]string{}
	var filename = dbFileName

	// load data from disk
	var data, _ = util.ReadJsonObjectDisk(filename)

	if len(data) > 0 {
		db = data
	}

	fmt.Println("Welcome to the simplest key-value memory database")
	for {

		fmt.Print("db>")
		fmt.Scan(&cmd)

		if cmd == "add" {

			var key string
			var value string

			fmt.Print("Enter key: ")
			fmt.Scan(&key)

			// check if key already exists
			if _, ok := db[key]; ok {
				fmt.Println("Key already exists")
			} else {
				fmt.Print("Enter value: ")
				fmt.Scan(&value)

				db[key] = value
				fmt.Println("Key added")
				util.WriteJsonDisk(filename, db)
			}

		} else if cmd == "update" {

			var key string
			var value string

			fmt.Print("Enter key: ")
			fmt.Scan(&key)

			// check if key already exists
			if _, ok := db[key]; ok {
				fmt.Print("Enter value: ")
				fmt.Scan(&value)

				db[key] = value
				fmt.Println("Key updated")
				util.WriteJsonDisk(filename, db)
			} else {
				fmt.Println("Key not exists")
			}

		} else if cmd == "get" {

			var key string
			fmt.Print("Enter key: ")
			fmt.Scan(&key)
			fmt.Println(db[key])

		} else if cmd == "delete" {

			var key string
			fmt.Print("Enter key: ")
			fmt.Scan(&key)
			// check if key exists
			if _, ok := db[key]; ok {
				delete(db, key)
				fmt.Println(key + " deleted")
				util.WriteJsonDisk(filename, db)
			} else {
				fmt.Println("Key not exists")
			}

		} else if cmd == "list" {
			for i, v := range db {
				fmt.Println(i, v)
			}
		} else if cmd == "exit" {
			fmt.Println("Bye")
			break
		} else {
			fmt.Println("Invalid command")
		}
	}
}
