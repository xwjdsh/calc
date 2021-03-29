package calc

import (
	"fmt"

	"github.com/xwjdsh/calc/operator"
	"github.com/xwjdsh/calc/variable"
)

type Formula struct {
	c *Calculator
}

func Compile(input string) (*Formula, error) {
	c := New().CompileMode()
	if _, err := c.eval(input, nil); err != nil {
		return nil, err
	}
	for {
		ok, err := c.compileCalculator.executeLastWithCondition(nil)
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
	}

	//fmt.Println("operators: ", c.compileCalculator.operatorStack.String())
	//fmt.Println("params: ", c.compileCalculator.paramStack.String())
	//fmt.Println("elements: ", c.compileCalculator.elements.String())
	return &Formula{c: c.compileCalculator}, nil
}

func (f *Formula) Eval(m map[string]interface{}) (interface{}, error) {
	nc := New()
	for _, ele := range f.c.elements.All() {
		if op, ok := ele.(operator.Operator); ok {
			preIsOperator := false
			if err := nc.addOperator(op, &preIsOperator); err != nil {
				return nil, err
			}
		} else {
			if v, ok := ele.(*variable.Var); ok {
				if value, ok := m[v.Name()]; ok {
					vv, ok := convertAndValidation(value)
					if !ok {
						return nil, fmt.Errorf("calc: unsupported param type %T", value)
					}
					nc.pushParam(vv)
				} else {
					return nil, fmt.Errorf("calc: unprovided variable: %s", v.Name())
				}
			} else {
				nc.pushParam(ele)
			}
		}
	}

	if err := nc.executeAll(); err != nil {
		return nil, err
	}

	result, ok := nc.paramStack.Top()
	if !ok {
		return nil, fmt.Errorf("calc: unable to parse expressions")
	}

	return result, nil
}
