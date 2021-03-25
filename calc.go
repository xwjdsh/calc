package calc

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"strconv"

	"github.com/xwjdsh/calc/operator"
	"github.com/xwjdsh/calc/stack"
)

var defaultCalculator = New()

type Calculator struct {
	paramStack    *stack.Stack
	operatorStack *stack.Stack
	opManager     *operator.Manager
}

func New() *Calculator {
	return &Calculator{
		paramStack:    stack.New(),
		operatorStack: stack.New(),
		opManager:     operator.NewManager(),
	}
}

func Eval(input string) (interface{}, error) {
	return defaultCalculator.Eval(input)
}

func (c *Calculator) Eval(input string) (interface{}, error) {
	return defaultCalculator.EvalWithVars(input, nil)
}

func EvalWithVars(str string, m map[string]float64) (interface{}, error) {
	return defaultCalculator.EvalWithVars(str, m)
}

func (c *Calculator) EvalWithVars(input string, m map[string]float64) (interface{}, error) {
	c.operatorStack.Clear()
	c.paramStack.Clear()

	var s scanner.Scanner
	src := []byte(input)
	fileSet := token.NewFileSet()
	s.Init(fileSet.AddFile("", fileSet.Base(), len(src)), src, nil, 0)

	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		switch {
		case tok == token.FLOAT || tok == token.INT:
			// is number
			f, err := strconv.ParseFloat(lit, 64)
			if err != nil {
				return nil, err
			}
			c.paramStack.Push(f)
		case tok == token.CHAR || tok == token.STRING:
			// is string
			// remove surrounding double or single quotes
			c.paramStack.Push(lit[1 : len(lit)-1])
		case c.opManager.Contains(tok.String()):
			// is operator
			op := c.opManager.Get(tok.String())
			switch op.Type() {
			case operator.General:
				if err := c.handleGeneralOperator(op); err != nil {
					return nil, err
				}
			case operator.Bracket:
				if err := c.handleBracketOperator(op); err != nil {
					return nil, err
				}
			}
		}
	}

	if err := c.executeAll(); err != nil {
		return nil, err
	}

	result, ok := c.paramStack.Top()
	if !ok {
		return nil, fmt.Errorf("calc: unable to parse expressions")
	}

	return result, nil
}

func (c *Calculator) handleGeneralOperator(op operator.Operator) error {
	for {
		preOp, ok := mustOperator(c.operatorStack.Top())
		if !ok {
			break
		}

		if op.Preference() > preOp.Preference() {
			break
		}

		preOp, _ = mustOperator(c.operatorStack.Pop())
		if err := c.execute(preOp); err != nil {
			return err
		}
	}

	c.operatorStack.Push(op)
	return nil
}

func (c *Calculator) handleBracketOperator(op operator.Operator) error {
	switch op.String() {
	case operator.LPAREN:
		c.operatorStack.Push(op)
	case operator.RPAREN:
		for {
			preOp, ok := mustOperator(c.operatorStack.Top())
			if !ok {
				break
			}
			if preOp.String() == operator.LPAREN {
				break
			}
			if preOp.Type() != operator.General {
				return errors.New("calc: unexpected operator type")
			}

			preOp, _ = mustOperator(c.operatorStack.Pop())
			if err := c.execute(preOp); err != nil {
				return err
			}
		}

		// pop `(`
		if _, ok := mustOperator(c.operatorStack.Pop()); !ok {
			return errors.New("calc: can not find matching parenthesis")
		}
	}

	return nil
}

func (c *Calculator) executeAll() error {
	for {
		op, ok := mustOperator(c.operatorStack.Pop())
		if !ok {
			return nil
		}

		if err := c.execute(op); err != nil {
			return err
		}
	}
}

func (c *Calculator) execute(op operator.Operator) error {
	count := op.ArgsCount()
	args := make([]interface{}, count)
	for i := count - 1; i >= 0; i-- {
		arg, ok := c.paramStack.Pop()
		if !ok {
			return fmt.Errorf("calc: no enough params for operator: %s", op.String())
		}

		args[i] = arg
	}

	result, err := op.Execute(args)
	if err != nil {
		return err
	}

	c.paramStack.Push(result)
	return nil
}

func mustOperator(r interface{}, ok bool) (operator.Operator, bool) {
	if !ok {
		return nil, false
	}

	return r.(operator.Operator), true
}
