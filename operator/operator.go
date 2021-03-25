package operator

import (
	"fmt"
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
)

var (
	_ Operator = new(generalOperator)
	_ Operator = new(bracketOperator)
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
				return i1 % i2, nil
			}
		}
	case COMMA:
		if (oks1 && oks2) || (okf1 && okf2) {
			return []interface{}{arg1, arg2}, nil
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

func (o *bracketOperator) ArgsCount() int {
	panic(fmt.Sprintf("calc/operator: access ArgsCount for bracket opreator: %s", o.code))
}

func (o *bracketOperator) Execute(args []interface{}) (interface{}, error) {
	panic(fmt.Sprintf("calc/operator: access Execute for bracket opreator: %s", o.code))
}
