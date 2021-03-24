package operator

type Manager struct {
	m map[string]Operator
}

func NewManager() *Manager {
	m := map[string]Operator{}
	for _, c := range []string{Add, Sub, Mul, Div, Mod} {
		m[c] = &GeneralOperator{c}
	}
	return &Manager{
		m: m,
	}
}

func (m *Manager) Contains(code string) bool {
	return m.m[code] != nil
}

func (m *Manager) Get(code string) Operator {
	return m.m[code]
}
