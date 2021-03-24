package calc

import (
	"fmt"
	"go/scanner"
	"go/token"
	"strconv"

	"github.com/xwjdsh/calc/operator"
	"github.com/xwjdsh/calc/stack"
)

var defaultCalculator = New()

type Calculator struct {
	floatStack    *stack.Stack
	operatorStack *stack.Stack
	opManager     *operator.Manager
}

func New() *Calculator {
	return &Calculator{
		floatStack:    stack.New(),
		operatorStack: stack.New(),
		opManager:     operator.NewManager(),
	}
}

func Eval(str string, m map[string]float64) (interface{}, error) {
	return defaultCalculator.Eval(str, m)
}

func (c *Calculator) Eval(input string, m map[string]float64) (interface{}, error) {
	var s scanner.Scanner
	src := []byte(input)
	fileSet := token.NewFileSet()
	s.Init(fileSet.AddFile("", fileSet.Base(), len(src)), src, nil, 0)

	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s\t%s\t%q\n", fileSet.Position(pos), tok, lit)

		switch {
		case tok == token.FLOAT || tok == token.INT:
			// is number
			f, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return 0, err
			}
			c.floatStack.Push(f)

		case c.opManager.Contains(tok.String()):
			// is operator
			if err := c.addOperatorToStack(c.opManager.Get(tok.String())); err != nil {
				return 0, err
			}
		}
	}

	fmt.Println("float stack:", c.floatStack.String())
	fmt.Println("operator stack:", c.operatorStack.String())
	return 0, nil
}

func (c *Calculator) addOperatorToStack(op operator.Operator) error {
	for {
		preOp, ok := mustOperator(c.operatorStack.Top())
		if !ok {
			break
		}

		if op.Preference() > preOp.Preference() {
			break
		}
		// TODO

	}

	c.operatorStack.Push(op)
	return nil
}

func (c *Calculator) isOperator(tok string) bool {
	return c.opManager.Get(tok) != nil
}
