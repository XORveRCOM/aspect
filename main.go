package main

import (
	"fmt"
	fpath "path/filepath"
	"time"

	"github.com/xorvercom/util/pkg/json"
)

var (
	DATAPATH string
)

func main() {
	fmt.Println()
	if abs, err := fpath.Abs("./out"); err != nil {
		fmt.Printf("DATAPATH error: %v", err)
		return
	} else {
		DATAPATH = abs

		j := json.NewElemObject()
		j.Put("key", json.NewElemString("value"))
		utcstr := time.Now().UTC().Format("20060102150405")
		json.SaveToJSONFile(fpath.Join(DATAPATH, "test."+utcstr+".json"), j, true)
	}
}
