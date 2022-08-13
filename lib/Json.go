package lib

import "encoding/json"

// parse json | json decode
// var input = `{
// 	"name": "John",
// 	"age": 30,
// 	"city": "New York"
// }`
func ParsedJSON(input []byte, data any) {

	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		panic(err)
	}
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
