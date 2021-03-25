package operator

// Manager manage all available operator.
type Manager struct {
	m map[string]Operator
}

// NewManager returns a new manager instance.
func NewManager() *Manager {
	m := map[string]Operator{}
	// register general type operators
	for _, c := range []string{ADD, SUB, MUL, QUO, REM, COMMA} {
		m[c] = newGeneralOperator(c)
	}
	// register bracket type operators
	for _, c := range []string{LPAREN, RPAREN} {
		m[c] = newBracketOperator(c)
	}

	return &Manager{
		m: m,
	}
}

// Contains returns the given code if available.
func (m *Manager) Contains(code string) bool {
	return m.m[code] != nil
}

// Get returns operator by special code, it will be nil if not found.
func (m *Manager) Get(code string) Operator {
	return m.m[code]
}
