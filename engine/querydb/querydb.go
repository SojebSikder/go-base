package querydb

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
	"github.com/sojebsikder/go-base/util"
	// "github.com/Jeffail/gabs"
)

// db file name
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
				// reader := bufio.NewReader(os.Stdin)
				// text, _ := reader.ReadString('\n')
				Precompile(string("create db [" + text + "]"))
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
				Precompile(string(text))
			}

		}
	}

}

// break whole text into statements using semiclon(;)
// this is the entry point of querydb
func Precompile(text string) {
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

	// parsing
	extractDoc := Parser(text, reBracket)
	extractProperty := Parser(text, reProperty)
	extractData := Parser(text, reData)
	tokens := Tokenize(text, " ")

	// process main oprations
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
		docName := extractDoc[0].(string)
		property := extractProperty
		data := extractData

		// exsting data in db file
		jsonDocs, _ := readJsonFile()

		marge := []any{}
		var concatMap map[string]string

		for j := 0; j < len(property); j++ {
			_property := property[j].(string)
			_data := data[j].(string)
			_map := map[string]string{_property: _data}

			for i := 0; i < len(jsonDocs); i++ {
				// fmt.Print(jsonDocs[i].(map[string]any)["user"])

				if value, ok := jsonDocs[i].(map[string]any)[docName]; ok {
					// TODO concat map
					concatMap = _map
					// value.(map[string]any)[docName] = _map
					value.(map[string]any)[docName] = concatMap
					marge = append(marge, value)
				}
			}

		}

		fmt.Println(marge)
		util.WriteDisk(dbFileName, marge)

		// for i := 0; i < len(jsonDocs); i++ {
		// 	// fmt.Print(jsonDocs[i].(map[string]any)["user"])

		// 	if value, ok := jsonDocs[i].(map[string]any)[docName]; ok {
		// 		for j := 0; j < len(property); j++ {
		// 			_property := property[j].(string)
		// 			_data := data[j].(string)
		// 			_map := map[string]string{_property: _data}

		// 			fmt.Println(property)

		// 			value.(map[string]any)[docName] = _map
		// 			marge = append(marge, value)
		// 		}
		// 	}
		// }
		// fmt.Println(marge)
		// writeData(dbFileName, marge)
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
		util.WriteDisk(dbFileName, data)
	} else {
		marge := append(jsonDocs, db)
		util.WriteDisk(dbFileName, marge)
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
