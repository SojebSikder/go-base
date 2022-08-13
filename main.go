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
)

func main() {
	arg := os.Args[1]
	fileName := os.Args[2]

	// query := "create [user]"
	// query := "add [user] 'sojeb' 'sikder'"

	if arg == "run" {
		_, err := os.Stat(fileName)

		if errors.Is(err, os.ErrNotExist) {
			fmt.Println("file does not exist")
		} else {

			content, err := ioutil.ReadFile(fileName)

			if err != nil {
				log.Fatal(err)
			}
			compile(string(content))
		}
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
	tokens := Tokenize(text)

	// crud oprations
	switch tokens[0] {
	case "create":
		// create db document
		docName := extractDoc[0]
		createDbDoc(docName)
	case "add":
		// add data to db document
		// fmt.Println(m["user"])
		dbData := readJsonFile()

		docName := extractDoc[0]
		data := extractData

		for i := 0; i < len(data); i++ {
			fmt.Println(data[i])
			dbData[extractDoc[0]] = data[i]
			writeDataToDoc("db.json", docName, dbData)
		}
	default:
		fmt.Println("Invalid commad")

	}
}

// get tokens from query
func Tokenize(text string) []string {
	keywords := strings.Split(text, " ")
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
		// fmt.Println(element)
		arr = append(arr, element)
	}
	return arr
}

// create db document
func createDbDoc(docName string) {
	var db = map[string]string{}
	db[docName] = ""
	// createDbfile()
	appendDataToDbfile("db.json", db)
}

// create file for database
func createDbfile() {
	emptyFile, err := os.Create("db.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(emptyFile)
	emptyFile.Close()
}

// read database file
func readDbfile() {
	file, err := os.Open("db.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func readJsonFile() map[string]string {
	file, _ := ioutil.ReadFile("db.json")
	m := map[string]string{}
	err := json.Unmarshal([]byte(file), &m)
	if err != nil {
		panic(err)
	}
	return m
}

// write data to database file
func appendDataToDbfile(filename string, content any) {
	file, _ := json.MarshalIndent(content, "", " ")
	_ = ioutil.WriteFile(filename, file, 0644)
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
