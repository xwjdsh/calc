package calc

import (
	"errors"
	"fmt"
	"go/scanner"
	"go/token"
	"strconv"
	"strings"

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
	return defaultCalculator.EvalWithVars(input, nil)
}

func (c *Calculator) Eval(input string) (interface{}, error) {
	return c.EvalWithVars(input, nil)
}

func EvalWithVars(str string, m map[string]float64) (interface{}, error) {
	return defaultCalculator.EvalWithVars(str, m)
}

func (c *Calculator) EvalWithVars(input string, m map[string]float64) (interface{}, error) {
	c.operatorStack.Clear()
	c.paramStack.Clear()

	var s scanner.Scanner
	src := []byte(strings.ToLower(input))
	fileSet := token.NewFileSet()
	s.Init(fileSet.AddFile("", fileSet.Base(), len(src)), src, nil, 0)

	var preIsOperator *bool
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		//fmt.Printf("%s\t%s\n", tok, lit)

		op, isOperator := c.getOperator(tok, lit)
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
		case isOperator:
			// is operator
			switch op.Type() {
			case operator.General:
				if err := c.handleGeneralOperator(op, preIsOperator); err != nil {
					return nil, err
				}
			case operator.Bracket:
				if err := c.handleBracketOperator(op); err != nil {
					return nil, err
				}
			case operator.Function:
				c.operatorStack.Push(op)
			}
		case tok == token.SEMICOLON:
		default:
			return nil, fmt.Errorf("calc: unsupported token: %s", lit)
		}

		preIsOperator = &isOperator
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

func (c *Calculator) handleGeneralOperator(op operator.Operator, preIsOperator *bool) error {
	// need handle `-` specially, convert to opposite number function
	if op.String() == operator.SUB && (preIsOperator == nil || *preIsOperator) {
		op = c.opManager.Get(operator.OPP)
	}

	for {
		ok, err := c.executeLastWithCondition(func(lastOp operator.Operator) (bool, error) {
			return op.Preference() <= lastOp.Preference(), nil
		})
		if err != nil {
			return err
		}
		if !ok {
			break
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
			ok, err := c.executeLastWithCondition(func(lastOp operator.Operator) (bool, error) {
				return lastOp.String() != operator.LPAREN, nil
			})
			if err != nil {
				return err
			}
			if !ok {
				break
			}
		}

		// pop `(`
		if _, ok := mustOperator(c.operatorStack.Pop()); !ok {
			return errors.New("calc: can not find matching parenthesis")
		}

		// calculation if pre operator is function type
		if _, err := c.executeLastWithCondition(func(lastOp operator.Operator) (bool, error) {
			return lastOp.Type() == operator.Function, nil
		}); err != nil {
			return err
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

func (c *Calculator) executeLastWithCondition(conditionFunc func(lastOp operator.Operator) (bool, error)) (bool, error) {
	preOp, ok := mustOperator(c.operatorStack.Top())
	if !ok {
		return false, nil
	}
	ok, err := conditionFunc(preOp)
	if !ok || err != nil {
		return ok, err
	}

	preOp, _ = mustOperator(c.operatorStack.Pop())
	if err := c.execute(preOp); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Calculator) getOperator(tok token.Token, lit string) (operator.Operator, bool) {
	if c.opManager.Contains(tok.String()) || c.opManager.Contains(lit) {
		op := c.opManager.Get(tok.String())
		if op == nil {
			op = c.opManager.Get(lit)
		}
		return op, true
	}

	return nil, false
}

func mustOperator(r interface{}, ok bool) (operator.Operator, bool) {
	if !ok {
		return nil, false
	}

	return r.(operator.Operator), true
}
