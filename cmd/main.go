package main

import (
	"fmt"

	"github.com/xorvercom/util/pkg/json"
)

func main() {
	fmt.Println()
	j := json.NewElemObject()
	j.Put("key", json.NewElemString(""))
	json.SaveToJSONFile("test.json", j, true)
}
