# calc

A simple command-line calculator written in Go.

## Usage

```go
package main

import (
	"fmt"

	"github.com/xwjdsh/calc"
)

func main() {
	r := calc.Must(calc.Eval("max(1,2,3)+min(4,5,6)*$a+$b", map[string]interface{}{"a": 7, "b": 8.1}))
	fmt.Println(r)
}
```

## Command-lint tool

```shell
> cd project_root_dir
> go build ./cmd/calc
>
> ./calc "((1*2-3)+5/4)"
0.25
> ./calc "((1*2-3)+5)/4"
1
> ./calc "1+sin(cos(5))*2"
1.5597467015370547
> ./calc "sum(1,2,3)+pow(2,3-1)"
10
> ./calc -m '{"a":1,"b":2}' '$a+$b'
3
> ./calc -m '{"a":1,"b":2}' '$a+$b'
3
> ./calc '"hello "*2+"world"'
hello hello world
>
```

### Supported Operators

##### General
`+`, `-`, `*`, `/`, `%`, `,`

##### Function
`sin`, `cos`, `tan`, `min`, `max`, `sum`, `pow`, `abs`, `opp`

##### Bracket
`(`, `)`

## Licence
[MIT License](https://github.com/xwjdsh/calc/blob/main/LICENSE)