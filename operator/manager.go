package operator

import "github.com/xwjdsh/calc/operator/token"

var opMap = map[token.T]Operator{}

func init() {
	// register general type operators
	for _, c := range []token.T{token.ADD, token.SUB, token.MUL, token.QUO, token.REM, token.COMMA} {
		opMap[c] = newGeneralOperator(c)
	}

	// register bracket type operators
	for _, c := range []token.T{token.LPAREN, token.RPAREN} {
		opMap[c] = newBracketOperator(c)
	}

	// register function type operators
	for _, c := range []token.T{token.SIN, token.COS, token.TAN, token.ABS, token.OPP, token.SUM, token.MAX, token.MIN, token.POW, token.REC} {
		opMap[c] = newFunctionOperator(c)
	}
}

// Manager manage all available operator.
type Manager struct {
	m map[token.T]Operator
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	return &Manager{
		m: opMap,
	}
}

func (m *Manager) SafeGet(t token.T) Operator {
	return m.m[t]
}

// Get returns operator by special code, it will be nil if not found.
func (m *Manager) Get(t token.T) (Operator, bool) {
	op, ok := m.m[t]
	return op, ok
}

// GetByString same as Get, but accept string.
func (m *Manager) GetByString(code string) (Operator, bool) {
	return m.Get(token.T(code))
}
