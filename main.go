// Author : sojebsikder<sojebsikder@gmail.com>
package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/sojebsikder/go-base/lib"
)

func main() {

	// app info
	appName := "go-base"
	version := "0.0.1"
	usage := "Welcome to go-base"
	description := "go-db is a simple database application"
	fmt.Printf("%s %s - %s\n", appName, version, usage)
	//

	if len(os.Args) < 2 {
		// run interactive mode
	} else {
		arg := os.Args[1]
		if arg == "version" {
			fmt.Println(appName + ": " + version)
		} else if arg == "help" {
			fmt.Println(description)

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
				lib.Precompile(string(content))
			}
		} else if arg == "mapdb" {
			lib.MapDB()
		} else if arg == "cli" {
			lib.Cli()
		} else {
			fmt.Println("Invalid command")
		}
	}

}
