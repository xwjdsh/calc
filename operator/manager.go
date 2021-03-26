package operator

// Manager manage all available operator.
type Manager struct {
	m map[token]Operator
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	m := map[token]Operator{}
	// register general type operators
	for _, c := range []token{ADD, SUB, MUL, QUO, REM, COMMA} {
		m[c] = newGeneralOperator(c)
	}

	// register bracket type operators
	for _, c := range []token{LPAREN, RPAREN} {
		m[c] = newBracketOperator(c)
	}

	// register function type operators
	for _, c := range []token{SIN, COS, TAN, ABS, OPP, SUM, MAX, MIN, POW} {
		m[c] = newFunctionOperator(c)
	}

	return &Manager{
		m: m,
	}
}

// Get returns operator by special code, it will be nil if not found.
func (m *Manager) Get(t token) (Operator, bool) {
	op, ok := m.m[t]
	return op, ok
}

// GetByString same as Get, but accept string.
func (m *Manager) GetByString(code string) (Operator, bool) {
	return m.Get(token(code))
}
