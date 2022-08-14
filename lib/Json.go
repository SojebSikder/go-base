package lib

import (
	"encoding/json"
	"errors"
)

// parse json | json decode
// string to json
// string -> json
//
//	var input = `{
//		"name": "John",
//		"age": 30,
//		"city": "New York"
//	}`
func ParsedJSON(input []byte, data any) error {

	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		// panic(err)
		return errors.New("something went wrong")
	}
	return nil
}

// string to json
// string -> json
func StringToJSON(data string) any {
	var any any
	ParsedJSON([]byte(data), &any)
	return any
}

// json encode
// stringify json
// json -> string
func Stringify(data any) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}
