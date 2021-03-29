package calc

import (
	"github.com/xwjdsh/calc/operator"
)

func (c *Calculator) getCompileCalculator() *Calculator {
	cc := c.compileCalculator
	if c.isCompileCalculator {
		cc = c
	}

	return cc
}

func (c *Calculator) clearElements() {
	c.operatorStack.Clear()
	c.paramStack.Clear()
	c.elements.Clear()
}

func (c *Calculator) pushOperator(ele interface{}) {
	c.operatorStack.Push(ele)
	c.elements.Push(ele)
}

func (c *Calculator) addFirstOperator(ele interface{}) {
	c.operatorStack.AddFirst(ele)
	c.elements.AddFirst(ele)
}

func (c *Calculator) pushParam(ele interface{}) {
	c.paramStack.Push(ele)
	c.elements.Push(ele)
}

func (c *Calculator) popOperator() (interface{}, bool) {
	r, ok := c.operatorStack.Pop()
	if ok {
		c.elements.Pop()
	}
	return r, ok
}

func (c *Calculator) popParam() (interface{}, bool) {
	r, ok := c.paramStack.Pop()
	if ok {
		c.elements.Pop()
	}
	return r, ok
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
