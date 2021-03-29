package calc

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

func BenchmarkEval(b *testing.B) {
	f, err := Compile(`max((1 + -2.12 + (2.12 - 3.0) + 10) * sin($column_123) / pow(2, 3), 1, 2)`)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println("formula: ", f.c.elements.String())

	for i := 0; i < b.N; i++ {
		_, _ = f.Eval(map[string]interface{}{"column_123": rand.Float64()})
	}
}
