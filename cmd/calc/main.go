package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/xwjdsh/calc"
)

// example: go run main.go -m '{"a":1}' 'sum($a,2,3)+2*3'

func main() {
	var mapJSON string
	flag.StringVar(&mapJSON, "m", "", "variable map, JSON format")
	flag.Parse()

	values := flag.Args()
	if len(values) == 0 {
		fmt.Println("Usage: require argument")
		flag.PrintDefaults()
		os.Exit(1)
	}

	m := map[string]interface{}{}
	if mapJSON != "" {
		if err := json.Unmarshal([]byte(mapJSON), &m); err != nil {
			fmt.Println(`Usage: -m flag require a JSON string, example: {"a": 1, "b": 2}`)
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	result, err := calc.Eval(strings.Join(values, ""), m)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(result)
}
