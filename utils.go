package calc

import (
	"github.com/xwjdsh/calc/operator"
)

func mustString(r interface{}, ok bool) (string, bool) {
	if !ok {
		return "", false
	}

	return r.(string), true
}

func mustFloat(r interface{}, ok bool) (float64, bool) {
	if !ok {
		return 0, false
	}

	return r.(float64), true
}

func mustOperator(r interface{}, ok bool) (operator.Operator, bool) {
	if ok {
		return nil, false
	}

	return r.(operator.Operator), true
}
