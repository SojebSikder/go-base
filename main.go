package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	// "github.com/Jeffail/gabs"
)

// db file name
var dbFileNameDir string = "db"
var dbFileName string

func main() {

	// app info
	appName := "go-db"
	version := "0.0.1"
	usage := "Welcome to go-db"
	fmt.Printf("%s %s - %s\n", appName, version, usage)
	//

	if len(os.Args) < 2 {
		// run interactive mode
		// fmt.Println("Welcome to simple db named as go-db")
	} else {
		arg := os.Args[1]
		if arg == "version" {
			fmt.Println(appName + ": " + version)
		} else if arg == "help" {
			fmt.Println("go-db is a simple database application")
		} else if arg == "run" {
			fileName := os.Args[2]
			_, err := os.Stat(fileName)

			if errors.Is(err, os.ErrNotExist) {
				fmt.Println("file does not exist")
			} else {

				content, err := ioutil.ReadFile(fileName)

				if err != nil {
					log.Fatal(err)
				}
				precompile(string(content))
			}
		} else if arg == "cli" {
			// first ask for database name
			if dbFileName == "" {
				fmt.Print("Select database: ")
				// reader := bufio.NewReader(os.Stdin)
				// dbName, _ := reader.Read() .ReadString('\n')
				var dbName string
				fmt.Scan(&dbName)
				dbFileName = dbFileNameDir + "/" + strings.Trim(dbName, "\n") + ".json"

				_, err := os.Stat(dbFileName)

				if os.IsNotExist(err) {
					// create new database file
					fmt.Print("db> ")
					reader := bufio.NewReader(os.Stdin)
					text, _ := reader.ReadString('\n')
					precompile(string(text))
					fmt.Println("Database selected: " + dbFileName)

				} else {
					fmt.Println("Database selected: " + dbFileName)
				}
			}

			for {
				fmt.Print("db> ")
				reader := bufio.NewReader(os.Stdin)
				text, _ := reader.ReadString('\n')
				precompile(string(text))
			}

		} else {
			fmt.Println("Invalid command")
		}
	}

}

// break whole text into statements
func precompile(text string) {
	tokens := Tokenize(text, ";")
	for i := 0; i < len(tokens); i++ {
		compile(string(tokens[i]))
	}
}

// compile query
func compile(text string) {
	// regex for brackets
	reBracket := regexp.MustCompile(`\[([^\[\]]*)\]`)
	// regex for single quotes
	re := regexp.MustCompile(`'([^']*)'`)

	extractDoc := Parser(text, reBracket)
	extractData := Parser(text, re)
	tokens := Tokenize(text, " ")

	// crud oprations
	switch tokens[0] {
	case "create":
		fmt.Print(tokens[1])
		if tokens[1] == "db" {
			// create db file
			fmt.Println("Creating database")
			dbName := extractDoc[0]
			createDbfile(dbFileNameDir + "/" + dbName + ".json")
		} else {

			// create db document
			docName := extractDoc
			for i := 0; i < len(docName); i++ {
				createDbDoc(docName[i])
			}
		}
	case "insert":
		// add data to db document
		dbData := readJsonFile()

		docName := extractDoc[0]
		data := extractData
		arr := []any{}

		for i := 0; i < len(data); i++ {
			dbData[extractDoc[0]] = data[i]
			arr = append(arr, dbData)
		}
		writeDataToDoc(dbFileName, docName, arr)
	default:
		fmt.Println("Invalid commad")
	}
}

// get tokens from query
func Tokenize(text string, deli string) []string {
	if deli == "" {
		deli = " "
	}
	keywords := strings.Split(text, deli)
	return keywords
}

// parse query for brackets, single quote
func Parser(text string, re *regexp.Regexp) []string {
	var arr []string

	submatchall := re.FindAllString(text, -1)
	for _, element := range submatchall {
		element = strings.Trim(element, "'")
		element = strings.Trim(element, "[")
		element = strings.Trim(element, "]")
		arr = append(arr, element)
	}
	return arr
}

// create db document
func createDbDoc(docName string) {
	var db = map[string]string{}
	db[docName] = ""
	appendDataToDbfile(dbFileName, db)
}

// // create file for database
func createDbfile(dbName string) {
	// create directory if not exists
	_, err := os.Stat(dbFileNameDir)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dbFileNameDir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}

	}

	emptyFile, err := os.Create(dbName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(emptyFile)
	emptyFile.Close()
}

// // read database file
// func readDbfile() {
// 	file, err := os.Open("db.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer file.Close()

// 	scanner := bufio.NewScanner(file)
// 	for scanner.Scan() {
// 		fmt.Println(scanner.Text())
// 	}
// }

func readJsonFile() map[string]string {
	file, _ := ioutil.ReadFile(dbFileName)
	m := map[string]string{}
	err := json.Unmarshal([]byte(file), &m)
	if err != nil {
		panic(err)
	}
	return m
}

// write data to database file
func appendDataToDbfile(filename string, data any) {
	// file, _ := json.MarshalIndent(content, "", " ")
	// _ = ioutil.WriteFile(filename, file, 0644)

	file, _ := json.Marshal(data)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", file)
}

func writeDataToDoc(filename string, docName string, data any) {
	file, _ := json.Marshal(data)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", file)
}

// CLI based db operations
func DB() {

	var cmd string
	var db = map[string]string{}

	fmt.Println("Welcome to the simplest key-value database")
	for {

		fmt.Println("Enter command")
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
