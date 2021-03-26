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

func Eval(str string, m map[string]interface{}) (interface{}, error) {
	return defaultCalculator.Eval(str, m)
}

func (c *Calculator) Eval(input string, m map[string]interface{}) (interface{}, error) {
	c.operatorStack.Clear()
	c.paramStack.Clear()

	var s scanner.Scanner
	src := []byte(strings.ToLower(input))
	fileSet := token.NewFileSet()
	s.Init(fileSet.AddFile("", fileSet.Base(), len(src)), src, nil, 0)

	if m == nil {
		m = map[string]interface{}{}
	}

	var (
		preIsOperator *bool
		parseVariable bool
	)
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		//fmt.Printf("%s\t%s\n", tok, lit)
		//fmt.Println("operator: ", c.operatorStack.String())
		//fmt.Println("params: ", c.paramStack.String())

		op, isOperator := c.getOperator(tok, lit)
		switch {
		case parseVariable:
			vn := lit
			if vn == "" {
				vn = tok.String()
			}

			v, ok := m[vn]
			if !ok {
				return nil, fmt.Errorf("calc: unknown variable: %s", vn)
			}

			if nv, ok := convertAndValidation(v); ok {
				c.paramStack.Push(nv)
			} else {
				return nil, fmt.Errorf("calc: unsupported variable type, name: %s, type: %T", vn, v)
			}
			parseVariable = false
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
		case lit == "$":
			parseVariable = true
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
	if op.Token() == operator.SUB && (preIsOperator == nil || *preIsOperator) {
		op = c.opManager.Get(operator.OPP)
	}

	for {
		ok, err := c.executeLastWithCondition(func(lastOp operator.ExecutableOperator) (bool, error) {
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
	switch op.Token() {
	case operator.LPAREN:
		c.operatorStack.Push(op)
	case operator.RPAREN:
		for {
			ok, err := c.executeLastWithCondition(nil)
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
		if _, err := c.executeLastWithCondition(func(lastOp operator.ExecutableOperator) (bool, error) {
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

		eop, ok := op.(operator.ExecutableOperator)
		if !ok {
			return fmt.Errorf("calc: unexecutable operator: %s", op.Token())
		}

		if err := c.execute(eop); err != nil {
			return err
		}
	}
}

func (c *Calculator) execute(op operator.ExecutableOperator) error {
	count := op.ArgsCount()
	args := make([]interface{}, count)
	for i := count - 1; i >= 0; i-- {
		arg, ok := c.paramStack.Pop()
		if !ok {
			return fmt.Errorf("calc: no enough params for operator: %s", op.Token())
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

func (c *Calculator) executeLastWithCondition(conditionFunc func(lastOp operator.ExecutableOperator) (bool, error)) (bool, error) {
	preOp, ok := mustOperator(c.operatorStack.Top())
	if !ok {
		return false, nil
	}
	eop, ok := preOp.(operator.ExecutableOperator)
	if !ok {
		return false, nil
	}

	if conditionFunc != nil {
		ok, err := conditionFunc(eop)
		if !ok || err != nil {
			return ok, err
		}
	}

	preOp, _ = mustOperator(c.operatorStack.Pop())
	if err := c.execute(eop); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Calculator) getOperator(tok token.Token, lit string) (operator.Operator, bool) {
	if c.opManager.Contains(tok.String()) || c.opManager.Contains(lit) {
		op := c.opManager.GetByString(tok.String())
		if op == nil {
			op = c.opManager.GetByString(lit)
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

func convertAndValidation(i interface{}) (interface{}, bool) {
	var r interface{}
	switch v := i.(type) {
	case string:
		r = v
	case float64:
		r = v
	case uint:
		r = float64(v)
	case uint8:
		r = float64(v)
	case uint16:
		r = float64(v)
	case uint32:
		r = float64(v)
	case uint64:
		r = float64(v)
	case int:
		r = float64(v)
	case int8:
		r = float64(v)
	case int16:
		r = float64(v)
	case int32:
		r = float64(v)
	case int64:
		r = float64(v)
	case float32:
		r = float64(v)
	default:
		return nil, false
	}

	return r, true
}
