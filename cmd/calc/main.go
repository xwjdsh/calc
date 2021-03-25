package main

import (
	"fmt"

	"github.com/xwjdsh/calc"
)

func main() {
	fmt.Println(calc.Eval("sum(1+$b*$a,$c)", map[string]interface{}{"a": 1, "b": 2, "c": 3}))
}
