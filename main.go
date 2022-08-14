package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/sojebsikder/go-base/lib"
	// "github.com/Jeffail/gabs"
)

// db file name
var dbFileNameDir string = "db"
var dbFileName string

func main() {

	// app info
	appName := "go-base"
	version := "0.0.1"
	usage := "Welcome to go-base"
	fmt.Printf("%s %s - %s\n", appName, version, usage)
	//

	if len(os.Args) < 2 {
		// run interactive mode
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
				content, err := os.ReadFile(fileName)

				if err != nil {
					log.Fatal(err)
				}
				precompile(string(content))
			}
		} else if arg == "cli" {
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
						// reader := bufio.NewReader(os.Stdin)
						// text, _ := reader.ReadString('\n')
						precompile(string("create db [" + text + "]"))
						fmt.Println("Database created: " + dbFileName)
					} else {
						fmt.Println("So you want to create later.")
					}

				} else {
					fmt.Println("Database selected: " + dbFileName)
					for {
						// reader := bufio.NewReader(os.Stdin)
						// text, _ := reader.ReadString('\n')
						fmt.Print("db> ")
						reader := bufio.NewReader(os.Stdin)
						text, _ := reader.ReadString('\n')
						// var text string
						// fmt.Scan(&text)
						precompile(string(text))
					}

				}
			}

		} else {
			fmt.Println("Invalid command")
		}
	}

}

// break whole text into statements using semiclon(;)
func precompile(text string) {
	tokens := Tokenize(text, ";")
	for i := 0; i < len(tokens); i++ {
		compile(string(tokens[i]))
	}
}

// compile query
func compile(text string) {
	// regex for third brackets []
	reBracket := regexp.MustCompile(`\[([^\[\]]*)\]`)
	// regex for second backet {}
	reProperty := regexp.MustCompile(`\{([^\{\}]*)\}`)
	// regex for single quotes ''
	reData := regexp.MustCompile(`'([^']*)'`)

	extractDoc := Parser(text, reBracket)
	extractProperty := Parser(text, reProperty)
	extractData := Parser(text, reData)
	tokens := Tokenize(text, " ")

	// crud oprations
	switch tokens[0] {
	// query for db
	case "set":
		if tokens[1] == "db" {
			// set database
			dbName := extractDoc[0].(string)
			dbFileName = dbFileNameDir + "/" + dbName + ".json"
		}
	case "drop":
		if tokens[1] == "db" {
			// drop db file
			dbName := extractDoc[0].(string)
			deleteDbfile(dbFileNameDir + "/" + dbName + ".json")

		}
	// query for doc
	case "create":
		if tokens[1] == "db" {
			// create db file
			dbName := extractDoc[0].(string)
			createDbfile(dbFileNameDir + "/" + dbName + ".json")
		} else if tokens[1] == "doc" {
			// create db document
			docName := extractDoc
			for i := 0; i < len(docName); i++ {
				createDbDoc(docName[i].(string))
			}
		}
	// query for data
	case "insert":
		// add data to db document
		// dbData, _ := readJsonFile()

		docName := extractDoc[0].(string)
		property := extractProperty
		data := extractData

		arr := []any{}
		// marge := append(property, data...)

		for i := 0; i < len(property); i++ {
			_property := property[i].(string)
			_data := data[i].(string)
			_map := map[string]string{_property: _data}
			arr = append(arr, _map)
		}

		// for i := 0; i < len(data); i++ {
		// 	dbData[0] = data[i]
		// 	arr = append(arr, dbData)
		// 	fmt.Println(marge)
		// }
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
func Parser(text string, re *regexp.Regexp) []any {
	var arr []any

	submatchall := re.FindAllString(text, -1)
	for _, element := range submatchall {
		element = strings.Trim(element, "'")
		element = strings.Trim(element, "[")
		element = strings.Trim(element, "]")
		element = strings.Trim(element, "{")
		element = strings.Trim(element, "}")
		arr = append(arr, element)
	}
	return arr
}

// insert data to document
func InsertDataToDoc(docName string) {

	var db = map[string]string{}
	db[docName] = ""
	// exsting data in db file
	jsonDocs, error := readJsonFile()

	if error != nil {
		// Unmarshall to slice
		var data []any
		input := `[{"` + docName + `":""}]`
		// Unmarshal into slice
		lib.ParsedJSON([]byte(input), &data)
		// write data to db file
		writeData(dbFileName, data)
	} else {
		marge := append(jsonDocs, db)
		writeData(dbFileName, marge)
	}
}

// create db document
func createDbDoc(docName string) {

	var db = map[string]string{}
	db[docName] = ""
	// exsting data in db file
	jsonDocs, error := readJsonFile()

	if error != nil {
		// Unmarshall to slice
		var data []any
		input := `[{"` + docName + `":{}}]`
		// Unmarshal into slice
		lib.ParsedJSON([]byte(input), &data)
		// write data to db file
		writeData(dbFileName, data)
	} else {
		marge := append(jsonDocs, db)
		writeData(dbFileName, marge)
	}
}

// create file for database
func createDbfile(dbName string) {
	// create directory if not exists
	_, err := os.Stat(dbFileNameDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dbFileNameDir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	emptyFile, ferr := os.Create(dbName)
	if ferr != nil {
		log.Fatal(ferr)
	}
	emptyFile.Close()
}

// delete file for database
func deleteDbfile(dbName string) {
	// delete directory if not exists
	_, err := os.Stat(dbFileNameDir)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dbFileNameDir, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	ferr := os.Remove(dbName)
	if ferr != nil {
		log.Fatal(ferr)
	}
}

// read db document from json file
func readJsonFile() ([]any, error) {
	file, _ := os.ReadFile(dbFileName)
	var data []any
	error := lib.ParsedJSON(file, &data)
	if error != nil {
		return data, errors.New("something went wrong")
	}
	return data, nil
}

// append data to database file
func appendDataToDbfile(filename string, data any) {
	file, _ := json.Marshal(data)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", file)
}

// write data without appending
func writeData(filename string, data any) {
	file, _ := json.Marshal(data)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", file)
}

func writeDataToDoc(filename string, docName string, data any) {
	appendDataToDbfile(filename, data)
}

// CLI based simple db operations for store in memory
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
