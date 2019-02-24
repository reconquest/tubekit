package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	params := parseParams(os.Args)

	{
		marshaledXXX, _ := json.MarshalIndent(params, "", "  ")
		fmt.Printf("params: %s\n", string(marshaledXXX))
	}
}
