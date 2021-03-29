// Package operator represents the operator in expression,
// include three types, general, bracket and function.

package operator

import (
	"fmt"
	"math"
	"strings"

	"github.com/xwjdsh/calc/operator/token"
	"github.com/xwjdsh/calc/variable"
)

// Type defines operator type.
type Type int

const (
	General Type = iota
	Bracket
	Function
)

var (
	_ Operator = new(generalOperator)
	_ Operator = new(bracketOperator)
	_ Operator = new(functionOperator)
)

// Operator abstract the methods to fetch operator basic info.
type Operator interface {
	// Token is the operator token.
	Token() token.T
	// Type operator type.
	Type() Type
	// Preference represent the operator priority, the bigger the value, the higher the priority.
	Preference() int
}

var (
	_ ExecutableOperator = new(generalOperator)
	_ ExecutableOperator = new(functionOperator)
)

// ExecutableOperator abstract the methods to execute operator.
type ExecutableOperator interface {
	Operator
	// ArgsCount returns the required arguments count.
	ArgsCount() int
	// Execute the operator handler
	Execute(args []interface{}) (interface{}, error)
}

type generalOperator struct {
	token.T
}

func newGeneralOperator(t token.T) *generalOperator {
	return &generalOperator{t}
}

func (o *generalOperator) Type() Type {
	return General
}

func (o *generalOperator) ArgsCount() int {
	return 2
}

func (o *generalOperator) Preference() int {
	switch o.T {
	case token.ADD, token.SUB:
		return 1
	case token.MUL, token.QUO, token.REM:
		return 2
	case token.COMMA:
		return -1
	}

	return 0
}

func (o *generalOperator) Execute(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("calc/operator: invalid param count for code: %s, expected: 2, actual: %d", o.T, len(args))
	}

	arg1, arg2 := args[0], args[1]
	vs1, oks1 := arg1.(string)
	vs2, oks2 := arg2.(string)

	vf1, okf1 := arg1.(float64)
	vf2, okf2 := arg2.(float64)

	veg1, okVeg1 := arg1.(variable.ExpGenerator)
	veg2, okVeg2 := arg2.(variable.ExpGenerator)

	commonGenVarExp := func() variable.ExpGenerator {
		if okVeg1 {
			return veg1.GenVarExp(o.T, arg2)
		}
		if okVeg2 {
			return veg2.GenVarExp(o.T, arg1)
		}

		return nil
	}

	switch o.T {
	case token.ADD:
		if oks1 && oks2 {
			return vs1 + vs2, nil
		}
		if okf1 && okf2 {
			return vf1 + vf2, nil
		}

		if arg1 == arg2 && okVeg1 {
			return veg1.GenVarExp(token.MUL, 2.0), nil
		}

		if r := commonGenVarExp(); r != nil {
			return r, nil
		}
	case token.MUL:
		if okf1 && okf2 {
			return vf1 * vf2, nil
		}
		if (oks1 && okf2) || (oks2 && okf1) {
			s, f := vs2, vf1
			if oks1 {
				s, f = vs1, vf2
			}
			i := int(f)
			if f == float64(i) {
				if i < 0 {
					return "", nil
				}

				return strings.Repeat(s, i), nil
			}
		}

		if okVeg1 && okf2 && vf2 == 1.0 {
			return veg1, nil
		}
		if okVeg2 && okf1 && vf1 == 1.0 {
			return veg2, nil
		}
		if r := commonGenVarExp(); r != nil {
			return r, nil
		}
	case token.SUB:
		if okf1 && okf2 {
			return vf1 - vf2, nil
		}

		if arg1 == arg2 && okVeg1 {
			return 0.0, nil
		}

		if okVeg1 {
			return veg1.GenVarExp(o.T, arg2), nil
		}
		if okVeg2 {
			return veg2.WithFunction(token.OPP, true).GenVarExp(token.ADD, arg1), nil
		}
	case token.QUO:
		if okf1 && okf2 {
			if vf2 == 0 {
				return nil, fmt.Errorf("calc/operator: division by zero")
			}
			return vf1 / vf2, nil
		}

		if okVeg1 {
			return veg1.GenVarExp(o.T, arg2), nil
		}
		if okVeg2 {
			return veg2.WithFunction(token.REC, true).GenVarExp(token.MUL, arg1), nil
		}

	case token.REM:
		if okf1 && okf2 {
			i1, i2 := int(vf1), int(vf2)
			if vf1 == float64(i1) && vf2 == float64(i2) {
				return float64(i1 % i2), nil
			}
		}
	case token.COMMA:
		if oks1 || okf1 {
			return []interface{}{arg1, arg2}, nil
		}

		vSli1, okSli1 := arg1.([]interface{})
		if okSli1 && len(vSli1) > 0 {
			return append(vSli1, arg2), nil
		}
		vSli2, okSli2 := arg2.([]interface{})
		if okSli2 && len(vSli2) > 0 {
			return append(vSli2, arg1), nil
		}

		if okVeg1 {
			return veg1.Group(arg2), nil
		}
		if okVeg2 {
			return veg2.Group(arg1), nil
		}
	}

	return nil, fmt.Errorf("calc/operator: invalid arguments for code: %s", o.T)
}

type bracketOperator struct {
	token.T
}

func newBracketOperator(t token.T) *bracketOperator {
	return &bracketOperator{t}
}

func (o *bracketOperator) Type() Type {
	return Bracket
}

func (o *bracketOperator) Preference() int {
	return -1
}

type functionOperator struct {
	token.T
}

func newFunctionOperator(t token.T) *functionOperator {
	return &functionOperator{t}
}

func (o *functionOperator) Type() Type {
	return Function
}

func (o *functionOperator) ArgsCount() int {
	return 1
}

func (o *functionOperator) Execute(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("calc/operator: invalid param count for code: %s, expected: 1, actual: %d", o.T, len(args))
	}

	arg1 := args[0]
	vf1, okf1 := arg1.(float64)
	vExp1, okExp1 := arg1.(variable.ExpGenerator)
	originalFunc := func(f float64) float64 { return f }
	floatOpMap := map[token.T]func(f float64) float64{
		token.SIN: math.Sin,
		token.COS: math.Cos,
		token.TAN: math.Tan,
		token.ABS: math.Abs,
		token.OPP: func(f float64) float64 { return -f },
		token.REC: func(f float64) float64 { return 1 / f },
		token.SUM: originalFunc,
		token.MAX: originalFunc,
		token.MIN: originalFunc,
	}

	if floatOpMap[o.T] != nil {
		if okf1 {
			return floatOpMap[o.T](vf1), nil
		}
		if okExp1 {
			return vExp1.WithFunction(o.T, false), nil
		}
	}

	vSli1, okSli1 := arg1.([]interface{})
	vEg1, okEg1 := arg1.(*variable.ExpGroup)
	switch o.T {
	case token.SUM:
		if okSli1 && len(vSli1) > 0 {
			if sum := 0.0; supposeFloatSlice(vSli1, func(_ int, f float64) { sum += f }) {
				return sum, nil
			}
		}
		if okEg1 && len(vEg1.Items) > 1 {
			if sum := 0.0; supposeFloatSlice(vEg1.Items, func(_ int, f float64) { sum += f }) {
				vEg1.Items = []interface{}{}
				if sum != 0 {
					vEg1.Items = append(vEg1.Items, sum)
				}

				return vEg1.WithFunction(o.T, false), nil
			}
		}
	case token.MAX, token.MIN:
		eff := func(f1, f2 float64) bool { return f1 > f2 }
		esf := func(s1, s2 string) bool { return s1 > s2 }
		if o.T == token.MIN {
			eff = func(f1, f2 float64) bool { return f1 < f2 }
			esf = func(s1, s2 string) bool { return s1 < s2 }
		}

		if okSli1 && len(vSli1) > 0 {
			if r := 0.0; supposeFloatSlice(vSli1, func(i int, f float64) {
				if i == 0 || eff(f, r) {
					r = f
				}
			}) {
				return r, nil
			}
		}

		if okEg1 && len(vEg1.Items) > 1 {
			if r := 0.0; supposeFloatSlice(vEg1.Items, func(i int, f float64) {
				if i == 0 || eff(f, r) {
					r = f
				}
			}) {
				vEg1.Items = []interface{}{r}
				return vEg1.WithFunction(o.T, false), nil
			}
		}
		if r := 0.0; supposeFloatSlice(vSli1, func(i int, f float64) {
			if i == 0 || eff(f, r) {
				r = f
			}
		}) {
			return r, nil
		}

		if r := ""; supposeStringSlice(vSli1, func(i int, s string) {
			if i == 0 || esf(s, r) {
				r = s
			}
		}) {
			return r, nil
		}
	case token.POW:
		if len(vSli1) == 2 {
			fs := make([]float64, 2)
			if supposeFloatSlice(vSli1, func(i int, f float64) {
				fs[i] = f
			}) {
				return math.Pow(fs[0], fs[1]), nil
			}
		}
	}

	return nil, fmt.Errorf("calc/operator: invalid arguments for code: %s", o.T)
}

func (o *functionOperator) Preference() int {
	if o.T == token.OPP {
		return 2
	}

	return 0
}

func supposeFloatSlice(vs []interface{}, f func(i int, f float64)) bool {
	for i, v := range vs {
		tmp, ok := v.(float64)
		if !ok {
			return false
		}
		f(i, tmp)
	}

	return true
}

func supposeStringSlice(vs []interface{}, f func(i int, s string)) bool {
	for i, v := range vs {
		tmp, ok := v.(string)
		if !ok {
			return false
		}
		f(i, tmp)
	}

	return true
}
