package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	// DB()
	// str1 := "this is a [sample] [[string]] with [SOME] special words"
	// query := "create [user]"
	query := "add [user] 'sojeb' 'sikder'"

	// regex for brackets
	reBracket := regexp.MustCompile(`\[([^\[\]]*)\]`)
	// regex for single quotes
	re := regexp.MustCompile(`'([^']*)'`)

	compile(query, re, reBracket)
}

func compile(text string, re *regexp.Regexp, reBracket *regexp.Regexp) {
	extractDoc := Parser(text, reBracket)
	extractData := Parser(text, re)
	tokens := Tokenize(text)

	// crud oprations
	switch tokens[0] {
	case "create":
		docName := extractDoc[0]
		// create db document
		createDbDoc(docName)
	case "add":
		docName := extractDoc[0]
		data := extractData

		for i := 0; i < len(data); i++ {
			fmt.Println(data[i])
			writeDataToDoc("db.json", docName, data[i])
		}
	default:
		fmt.Println("Invalid commadn")

	}
}

func Tokenize(text string) []string {
	keywords := strings.Split(text, " ")
	// for _, keyword := range keywords {
	// 	if re.MatchString(keyword) {
	// 		fmt.Println(re.FindStringSubmatch(keyword))
	// 	}
	// }
	return keywords
}

func Parser(text string, re *regexp.Regexp) []string {
	var arr []string

	// fmt.Printf("Pattern: %v\n", re.String())      // print pattern
	// fmt.Println("Matched:", re.MatchString(text)) // true

	// fmt.Println("\nText between square brackets:")
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
	db[docName] = "{" + docName + ":{}}"
	// createDbfile()
	appendDataToDbfile("db.json", db[docName])
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

// write data to database file
func appendDataToDbfile(filename string, content string) {
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

	// _ = ioutil.WriteFile(filename, file, 0644)
}

// cli based db operations
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
