package operator

import (
	"fmt"
	"strings"
)

type Type int

const (
	General Type = iota
	SingleParamFunction
	MultiParamFunction
)

const (
	Add = "+"
	Sub = "-"
	Mul = "*"
	Div = "/"
	Mod = "%"
)

type Operator interface {
	Code() string
	Type() Type
	Preference() int
	Execute(args []interface{}) (interface{}, error)
}

type GeneralOperator struct {
	code string
}

func (o *GeneralOperator) Code() string {
	return o.code
}

func (o *GeneralOperator) Type() Type {
	return General
}

func (o *GeneralOperator) Preference() int {
	switch o.code {
	case Add, Sub:
		return 1
	case Mul, Div, Mod:
		return 2
	}

	return 0
}

func (o *GeneralOperator) Execute(args []interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("calc/operator: invalid param count for code: %s, expected: 2, actual: %d", o.code, len(args))
	}

	arg1, arg2 := args[0], args[1]
	vs1, oks1 := arg1.(string)
	vs2, oks2 := arg2.(string)

	vf1, okf1 := arg1.(float64)
	vf2, okf2 := arg2.(float64)

	switch o.code {
	case Add:
		if oks1 && oks2 {
			return vs1 + vs2, nil
		}
		if okf1 && okf2 {
			return vf1 + vf2, nil
		}
	case Mul:
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
	case Sub:
		if okf1 && okf2 {
			return vf1 - vf2, nil
		}
	case Div:
		if okf1 && okf2 {
			if vf2 == 0 {
				return nil, fmt.Errorf("calc/operator: division by zero")
			}
			return vf1 / vf2, nil
		}
	case Mod:
		if okf1 && okf2 {
			i1, i2 := int(vf1), int(vf2)
			if vf1 == float64(i1) && vf2 == float64(i2) {
				return i1 % i2, nil
			}
		}
	}

	return nil, fmt.Errorf("calc/operator: invalid arguments for code: %s", o.code)
}
