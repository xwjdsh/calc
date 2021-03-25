package operator

import (
	"fmt"
	"math"
	"strings"
)

// Type defines operator type.
type Type int

const (
	General Type = iota
	Bracket
	Function
)

const (
	// general type
	ADD   = "+"
	SUB   = "-"
	MUL   = "*"
	QUO   = "/"
	REM   = "%"
	COMMA = ","

	// bracket type
	LPAREN = "("
	RPAREN = ")"

	// function type
	SIN = "sin"
	COS = "cos"
	TAN = "tan"
	ABS = "abs"
	OPP = "opp" // opposite number
	SUM = "sum"
	MAX = "max"
	MIN = "min"
	POW = "pow"
)

var (
	_ Operator = new(generalOperator)
	_ Operator = new(bracketOperator)
	_ Operator = new(functionOperator)
)

// Operator abstract the function of operators.
type Operator interface {
	// String is the operator token.
	String() string
	// Type operator type.
	Type() Type
	// ArgsCount returns the required arguments count.
	ArgsCount() int
	// Preference represent the operator priority, the bigger the value, the higher the priority.
	Preference() int
	// Execute the operator handler
	Execute(args []interface{}) (interface{}, error)
}

type generalOperator struct {
	code string
}

func newGeneralOperator(c string) *generalOperator {
	return &generalOperator{code: c}
}

func (o *generalOperator) String() string {
	return o.code
}

func (o *generalOperator) Type() Type {
	return General
}

func (o *generalOperator) ArgsCount() int {
	return 2
}

func (o *generalOperator) Preference() int {
	switch o.code {
	case ADD, SUB:
		return 1
	case MUL, QUO, REM:
		return 2
	}

	return 0
}

func (o *generalOperator) Execute(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("calc/operator: invalid param count for code: %s, expected: 2, actual: %d", o.code, len(args))
	}

	arg1, arg2 := args[0], args[1]
	vs1, oks1 := arg1.(string)
	vs2, oks2 := arg2.(string)

	vf1, okf1 := arg1.(float64)
	vf2, okf2 := arg2.(float64)

	switch o.code {
	case ADD:
		if oks1 && oks2 {
			return vs1 + vs2, nil
		}
		if okf1 && okf2 {
			return vf1 + vf2, nil
		}
	case MUL:
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
	case SUB:
		if okf1 && okf2 {
			return vf1 - vf2, nil
		}
	case QUO:
		if okf1 && okf2 {
			if vf2 == 0 {
				return nil, fmt.Errorf("calc/operator: division by zero")
			}
			return vf1 / vf2, nil
		}
	case REM:
		if okf1 && okf2 {
			i1, i2 := int(vf1), int(vf2)
			if vf1 == float64(i1) && vf2 == float64(i2) {
				return float64(i1 % i2), nil
			}
		}
	case COMMA:
		if oks1 || okf1 {
			return []interface{}{arg1, arg2}, nil
		}
		vSli1, okSli1 := arg1.([]interface{})
		if okSli1 && len(vSli1) > 0 {
			return append(vSli1, arg2), nil
		}
	}

	return nil, fmt.Errorf("calc/operator: invalid arguments for code: %s", o.code)
}

type bracketOperator struct {
	*generalOperator
}

func newBracketOperator(c string) *bracketOperator {
	return &bracketOperator{newGeneralOperator(c)}
}

func (o *bracketOperator) Type() Type {
	return Bracket
}

func (o *bracketOperator) Preference() int {
	return -1
}

func (o *bracketOperator) ArgsCount() int {
	panic(fmt.Sprintf("calc/operator: access ArgsCount for bracket opreator: %s", o.code))
}

func (o *bracketOperator) Execute([]interface{}) (interface{}, error) {
	panic(fmt.Sprintf("calc/operator: access Execute for bracket opreator: %s", o.code))
}

type functionOperator struct {
	*generalOperator
}

func newFunctionOperator(c string) *functionOperator {
	return &functionOperator{newGeneralOperator(c)}
}

func (o *functionOperator) Type() Type {
	return Function
}

func (o *functionOperator) ArgsCount() int {
	return 1
}

func (o *functionOperator) Execute(args []interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("calc/operator: invalid param count for code: %s, expected: 1, actual: %d", o.code, len(args))
	}

	arg1 := args[0]
	vf1, okf1 := arg1.(float64)
	originalFunc := func(f float64) float64 { return f }
	floatOpMap := map[string]func(f float64) float64{
		SIN: math.Sin,
		COS: math.Cos,
		TAN: math.Tan,
		ABS: math.Abs,
		OPP: func(f float64) float64 { return -f },
		SUM: originalFunc,
		MAX: originalFunc,
		MIN: originalFunc,
	}

	if okf1 && floatOpMap[o.code] != nil {
		return floatOpMap[o.code](vf1), nil
	}

	vSli1, okSli1 := arg1.([]interface{})
	if okSli1 && len(vSli1) > 0 {
		switch o.code {
		case SUM:
			if sum := 0.0; supposeFloatSlice(vSli1, func(_ int, f float64) { sum += f }) {
				return sum, nil
			}
		case MAX, MIN:
			eff := func(f1, f2 float64) bool { return f1 > f2 }
			esf := func(s1, s2 string) bool { return s1 > s2 }
			if o.code == MIN {
				eff = func(f1, f2 float64) bool { return f1 < f2 }
				esf = func(s1, s2 string) bool { return s1 < s2 }
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
		case POW:
			if len(vSli1) == 2 {
				fs := make([]float64, 2)
				if supposeFloatSlice(vSli1, func(i int, f float64) {
					fs[i] = f
				}) {
					return math.Pow(fs[0], fs[1]), nil
				}
			}
		}
	}

	return nil, fmt.Errorf("calc/operator: invalid arguments for code: %s", o.code)
}

func (o *functionOperator) Preference() int {
	if o.code == OPP {
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
