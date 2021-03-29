// Package calc implements the logic of expression calculation by Eval or obj.Eval.

package calc

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/xwjdsh/calc/operator"
	"github.com/xwjdsh/calc/operator/token"
	"github.com/xwjdsh/calc/stack"
	"github.com/xwjdsh/calc/variable"
)

var defaultCalculator = New()

// Calculator wraps the logic of executing expressions.
type Calculator struct {
	elements      *stack.Stack
	paramStack    *stack.Stack
	operatorStack *stack.Stack
	opManager     *operator.Manager
	vm            *variable.Manager

	compileMode         bool
	compileExp          *variable.Exp
	compileCalculator   *Calculator
	isCompileCalculator bool
}

// New returns a new Calculator instance.
func New() *Calculator {
	return &Calculator{
		elements:      stack.New(),
		paramStack:    stack.New(),
		operatorStack: stack.New(),
		opManager:     operator.NewManager(),
		vm:            variable.NewManager(),
	}
}

func (c *Calculator) CompileMode() *Calculator {
	c.compileMode = true
	c.compileCalculator = New()
	c.compileCalculator.isCompileCalculator = true
	return c
}

// Must cause panic if err is not nil.
func Must(r interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return r
}

// Eval calculates given expressions.
func Eval(str string, m map[string]interface{}) (interface{}, error) {
	return defaultCalculator.Eval(str, m)
}

func (c *Calculator) Eval(input string, m map[string]interface{}) (interface{}, error) {
	return c.eval(input, m)
}

func (c *Calculator) eval(input string, m map[string]interface{}) (interface{}, error) {
	c.clearElements()

	var s scanner.Scanner
	s.Init(strings.NewReader(strings.ToLower(input)))
	// ignore scanner error message
	s.Error = func(s *scanner.Scanner, msg string) {}

	if m == nil {
		m = map[string]interface{}{}
	}

	var (
		preIsOperator *bool
		parseVariable bool
	)
	for {
		tok := s.Scan()
		text := s.TokenText()
		if tok == scanner.EOF {
			break
		}

		//fmt.Printf("%s\t%s\n", scanner.TokenString(tok), text)
		//fmt.Println("operator: ", c.operatorStack.String())
		//fmt.Println("params: ", c.paramStack.String())

		op, isOperator := c.opManager.GetByString(text)
		switch {
		case parseVariable:
			parseVariable = false
			if c.compileMode {
				c.pushParam(c.vm.Get(text))
			} else {
				v, ok := m[text]
				if !ok {
					return nil, fmt.Errorf("calc: unknown variable: %s", text)
				}

				if nv, ok := convertAndValidation(v); ok {
					c.pushParam(nv)
				} else {
					return nil, fmt.Errorf("calc: unsupported variable type, name: %s, type: %T", text, v)
				}
			}
		case tok == scanner.Float || tok == scanner.Int:
			// is number
			f, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return nil, err
			}
			c.pushParam(f)
		case tok == scanner.Char || tok == scanner.String:
			// is string
			// remove surrounding double or single quotes
			c.pushParam(text[1 : len(text)-1])
		case isOperator:
			if err := c.addOperator(op, preIsOperator); err != nil {
				return nil, err
			}
		case text == "$":
			parseVariable = true
		default:
			return nil, fmt.Errorf("calc: unsupported token: '%s'", text)
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

func (c *Calculator) addOperator(op operator.Operator, preIsOperator *bool) error {
	// is operator
	switch op.Type() {
	case operator.General:
		if c.isCompileCalculator {
			c.pushOperator(op)
			return nil
		}
		if err := c.handleGeneralOperator(op, preIsOperator); err != nil {
			return err
		}
	case operator.Bracket:
		if err := c.handleBracketOperator(op); err != nil {
			return err
		}
	case operator.Function:
		c.pushOperator(op)
	}

	return nil
}

func (c *Calculator) handleGeneralOperator(op operator.Operator, preIsOperator *bool) error {
	// need handle `-` specially, convert to opposite number function
	// preIsOperator == nil, e.g. `-1+2`
	// *preIsOperator is true, e.g. `2+ -1` or `(-1-2)`
	if op.Token() == token.SUB {
		if preIsOperator == nil {
			// `-1` -> `0-1`
			c.pushParam(0.0)
		} else if *preIsOperator {
			// `1--1` -> `1-opp(1)`
			op, _ = c.opManager.Get(token.OPP)
		}
	}

	for {
		ok, err := c.executeLastWithCondition(func(preOp operator.Operator) (bool, error) {
			return op.Preference() <= preOp.Preference(), nil
		})
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}

	c.pushOperator(op)
	return nil
}

func (c *Calculator) handleBracketOperator(op operator.Operator) error {
	switch op.Token() {
	case token.LPAREN:
		c.pushOperator(op)
	case token.RPAREN:

		if c.isCompileCalculator {
			c.pushOperator(op)
		} else {
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
			if _, ok := c.popOperator(); !ok {
				return errors.New("calc: can not find matching parenthesis")
			}
		}

		if c.compileMode {
			v, ok := c.paramStack.Top()
			if ok {
				cc := c.getCompileCalculator()
				if exp, ok := v.(*variable.Exp); ok {
					exp.WithBracket()
					if exp == c.compileExp {
						cc.addFirstOperator(c.opManager.SafeGet(token.LPAREN))
						cc.pushOperator(c.opManager.SafeGet(token.RPAREN))
					}
				}
			}
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
		ok, err := c.executeLastWithCondition(func(lastOp operator.Operator) (b bool, err error) {
			eop, ok := lastOp.(operator.ExecutableOperator)
			if !ok {
				return false, fmt.Errorf("calc: unexecutable operator: %s", eop.Token())
			}
			return true, nil
		})
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
	}
}

func (c *Calculator) execute(op operator.ExecutableOperator, args []interface{}) error {
	result, err := op.Execute(args)
	if err != nil {
		return err
	}
	if c.compileMode || c.isCompileCalculator {
		cc := c.getCompileCalculator()
		if ve, ok := result.(*variable.Exp); ok { //&& (c.compileExp == nil || c.compileExp == ve.From()) {
			if c.compileExp == nil || c.compileExp == ve.From() {
				//fmt.Println("-> before elements", cc.elements.String())
				//fmt.Println("-> firstTokenAndParams", ve.FirstTokenAndParams)
				//fmt.Println("-> tokenAndParams", ve.TokenAndParams)
				c.compileExp = ve
				for _, p := range ve.TokenAndParams {
					if tok, ok := p.(token.T); ok {
						preIsOperator := false
						if err := cc.addOperator(cc.opManager.SafeGet(tok), &preIsOperator); err != nil {
							return err
						}
						continue
					}

					cc.pushParam(p)
				}

				for i := len(ve.FirstTokenAndParams) - 1; i >= 0; i-- {
					p := ve.FirstTokenAndParams[i]
					if tok, ok := p.(token.T); ok {
						cc.addFirstOperator(cc.opManager.SafeGet(tok))
					}
				}

				if err := cc.executeLastForCompile(); err != nil {
					return err
				}
				//fmt.Println("-> after elements", cc.elements.String())
				//fmt.Println()
			} else {
				ve.TokenAndParams = append(ve.FirstTokenAndParams, ve.TokenAndParams...)
				ve.FirstTokenAndParams = []interface{}{}
			}
		} else if c.isCompileCalculator {
			cc.pushParam(result)
		}
	}

	if !c.isCompileCalculator {
		c.pushParam(result)
	}
	return nil
}

func (c *Calculator) getExecuteParams(eop operator.ExecutableOperator) ([]interface{}, error) {
	count := eop.ArgsCount()
	args := make([]interface{}, count)

	params := c.paramStack.All()
	paramIndex := len(params) - 1

	for i := count - 1; i >= 0; i-- {
		//arg, ok := c.paramStack.Pop()
		//if !ok {
		//	return nil, fmt.Errorf("calc: no enough params for operator: %s", eop.Token())
		//}
		if paramIndex < 0 {
			return nil, fmt.Errorf("calc: no enough params for operator: %s", eop.Token())
		}
		args[i] = params[paramIndex]
		paramIndex--
	}

	if c.isCompileCalculator {
		if _, ok := args[0].(*variable.Var); ok {
			return nil, nil
		}
	}

	for i := count - 1; i >= 0; i-- {
		c.popParam()
	}

	return args, nil
}

func (c *Calculator) executeLastWithCondition(conditionFunc func(lastOp operator.Operator) (bool, error)) (bool, error) {
	preOp, ok := mustOperator(c.operatorStack.Top())
	if !ok {
		return false, nil
	}
	if conditionFunc != nil {
		ok, err := conditionFunc(preOp)
		if !ok || err != nil {
			return ok, err
		}
	}

	eop, ok := preOp.(operator.ExecutableOperator)
	if !ok {
		return false, nil
	}

	args, err := c.getExecuteParams(eop)
	if err != nil {
		return false, err
	}
	if args == nil {
		return false, nil
	}

	c.popOperator()
	if err := c.execute(eop, args); err != nil {
		return false, err
	}

	return true, nil
}

func (c *Calculator) executeLastForCompile() error {
	if !c.isCompileCalculator {
		return fmt.Errorf("calc: wrong execution for executeLastForCompile")
	}

	for {
		ok, err := c.executeLastWithCondition(func(preOp operator.Operator) (bool, error) {
			operators := c.operatorStack.All()
			count := len(operators)
			if count > 1 {
				lastOne, lastTwo := operators[count-1], operators[count-2]
				eop1, okEop1 := lastOne.(operator.ExecutableOperator)
				eop2, okEop2 := lastTwo.(operator.ExecutableOperator)
				return okEop1 && okEop2 && eop2.Preference() <= eop1.Preference(), nil
			}

			return false, nil
		})
		if err != nil {
			return err
		}
		if !ok {
			break
		}
	}

	return nil
}
