package calc

import (
	"math"
	"reflect"
	"testing"
)

func TestEval(t *testing.T) {
	cases := []struct {
		expressions string
		expected    interface{}
		expectErr   bool
	}{
		// general
		{expressions: "1+2", expected: 3.0},
		{expressions: "-1+2", expected: 1.0},
		{expressions: "10%3/2", expected: 0.5},
		{expressions: "1,2", expected: []interface{}{1.0, 2.0}},
		{expressions: "1/3", expected: 1.0 / 3},
		{expressions: "1+2*3+4", expected: 11.0},
		{expressions: "1/2*3", expected: 1.5},
		{expressions: `'hello' +  " " + 'world'`, expected: "hello world"},
		{expressions: `'hello ' * 2 + "world"`, expected: "hello hello world"},

		// bracket
		{expressions: "(1+2)*3+4", expected: 13.0},
		{expressions: "(1*(2+1))*((3+1)*4)", expected: 48.0},

		// function
		{expressions: "aBs(-1)", expected: 1.0},
		{expressions: "3*sin(3+2*3)+3*4", expected: 3*math.Sin(3+2*3) + 3*4},
		{expressions: "1+abs(-5)*2", expected: 11.0},
		{expressions: "1+opp(5)*2", expected: -9.0},
		{expressions: "sum(-5)+2*sum(1,2,3)", expected: 7.0},
		{expressions: "max(-5,1,3,8)+2", expected: 10.0},
		{expressions: "min(-5,1,3,8)+2", expected: -3.0},
		{expressions: "max('a','n','f')*2", expected: "nn"},
		{expressions: "min('a','n','f')*2", expected: "aa"},
		{expressions: "sum(max(-5,1,3,8),min(-5,1,3,8))", expected: 3.0},

		// error conditions
		{expressions: "2/0", expectErr: true},
		{expressions: `"test"*3.3`, expectErr: true},
		{expressions: "3%1.1", expectErr: true},
		{expressions: `"test"-3`, expectErr: true},
		{expressions: `"test"/3`, expectErr: true},
	}

	for _, c := range cases {
		result, err := Eval(c.expressions)
		if c.expectErr && err == nil {
			t.Errorf("expect error, got nil, expressions: %s", c.expressions)
		}

		if !c.expectErr && err != nil {
			t.Errorf("expect no error, got %v, expressions: %s", err, c.expressions)
		}

		if !reflect.DeepEqual(c.expected, result) {
			t.Errorf("expected: %v, got: %v, expressions: %s", c.expected, result, c.expressions)
		}
	}
}
