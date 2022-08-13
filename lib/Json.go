package lib

import "encoding/json"

// parse json
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
