package util

import (
	"fmt"
	"os"

	"github.com/sojebsikder/go-base/lib"
)

// write to disk and remove previous data
func WriteDisk(filename string, data any) {
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
