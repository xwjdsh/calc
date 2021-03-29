package operator

import (
	"math"
	"reflect"
	"testing"

	"github.com/xwjdsh/calc/operator/token"
)

func TestGeneralOperator(t *testing.T) {
	cases := []struct {
		code      token.T
		args      []interface{}
		expected  interface{}
		expectErr bool
	}{
		{code: token.ADD, args: []interface{}{1.0, 2.0}, expected: 3.0},
		{code: token.SUB, args: []interface{}{1.0, 2.0}, expected: -1.0},
		{code: token.MUL, args: []interface{}{1.0, 2.0}, expected: 2.0},
		{code: token.QUO, args: []interface{}{1.0, 2.0}, expected: 0.5},
		{code: token.REM, args: []interface{}{10.0, 3.0}, expected: 1.0},
		{code: token.COMMA, args: []interface{}{1.0, 2.0}, expected: []interface{}{1.0, 2.0}},
		{code: token.COMMA, args: []interface{}{[]interface{}{1.0, 2.0}, 3.0}, expected: []interface{}{1.0, 2.0, 3.0}},
		{code: token.QUO, args: []interface{}{1.0, 0.0}, expectErr: true},
		{code: token.ADD, args: []interface{}{1.0, 0.0, 1.0}, expectErr: true},
		{code: token.REM, args: []interface{}{10.1, 3.0}, expectErr: true},
	}

	for _, c := range cases {
		result, err := newGeneralOperator(c.code).Execute(c.args)
		if c.expectErr && err == nil {
			t.Errorf("expect error, got nil")
		}

		if !c.expectErr && err != nil {
			t.Errorf("expect no error, got %v", err)
		}

		if !reflect.DeepEqual(c.expected, result) {
			t.Errorf("expected: %v, got: %v", c.expected, result)
		}
	}
}

func TestFunctionOperator(t *testing.T) {
	cases := []struct {
		code      token.T
		args      []interface{}
		expected  interface{}
		expectErr bool
	}{
		{code: token.SIN, args: []interface{}{10.0}, expected: math.Sin(10)},
		{code: token.COS, args: []interface{}{10.0}, expected: math.Cos(10)},
		{code: token.TAN, args: []interface{}{10.0}, expected: math.Tan(10)},
		{code: token.ABS, args: []interface{}{-10.0}, expected: 10.0},
		{code: token.OPP, args: []interface{}{10.0}, expected: -10.0},
		{code: token.SUM, args: []interface{}{[]interface{}{1.0, 2.0, 3.0}}, expected: 6.0},
		{code: token.SUM, args: []interface{}{10.0}, expected: 10.0},
		{code: token.MAX, args: []interface{}{[]interface{}{1.0, 2.0, 3.0}}, expected: 3.0},
		{code: token.MIN, args: []interface{}{[]interface{}{1.0, 2.0, 3.0}}, expected: 1.0},
		{code: token.POW, args: []interface{}{[]interface{}{2.0, 2.0}}, expected: 4.0},
		{code: token.COS, args: []interface{}{10.0, 2.0}, expectErr: true},
	}

	for _, c := range cases {
		result, err := newFunctionOperator(c.code).Execute(c.args)
		if c.expectErr && err == nil {
			t.Errorf("expect error, got nil")
		}

		if !c.expectErr && err != nil {
			t.Errorf("expect no error, got %v", err)
		}

		if !reflect.DeepEqual(c.expected, result) {
			t.Errorf("expected: %v, got: %v", c.expected, result)
		}
	}
}
