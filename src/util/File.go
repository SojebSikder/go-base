package util

import (
	// "errors"
	"fmt"
	"os"

	"github.com/sojebsikder/go-base/src/lib"
)

// write to disk and remove previous data
func WriteJsonDisk(filename string, data any) {
	// file, _ := json.Marshal(data)
	file := lib.Stringify(data)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_CREATE, 0660)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	defer f.Close()

	fmt.Fprintf(f, "%s\n", file)
}

// read data from disk
func ReadJsonDisk(filename string) ([]any, error) {
	file, _ := os.ReadFile(filename)
	var data []any
	error := lib.ParsedJSON(file, &data)
	if error != nil {
		// return data, errors.New("something went wrong")
		return data, error
	}
	return data, nil
}

// read data from disk
func ReadDisk(filename string) (any, error) {
	file, error := os.ReadFile(filename)
	var data = string(file)
	if error != nil {
		return data, error
	}
	return data, nil
}
