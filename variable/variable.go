package variable

import (
	"github.com/xwjdsh/calc/operator/token"
)

var (
	_ ExpGenerator = new(Var)
	_ ExpGenerator = new(Exp)
	_ ExpGenerator = new(ExpGroup)
)

type ExpGenerator interface {
	GenVarExp(ps ...interface{}) ExpGenerator
	WithFunction(t token.T, addBracket bool) *Exp
	Group(ps ...interface{}) *ExpGroup
}

type Manager struct {
	m map[string]*Var
}

func NewManager() *Manager {
	return &Manager{
		m: map[string]*Var{},
	}
}

func (m *Manager) Get(name string) *Var {
	v, ok := m.m[name]
	if !ok {
		v = newVar(name)
		m.m[name] = v
	}

	return v
}

type Var struct {
	name string
}

func newVar(name string) *Var {
	v := &Var{
		name: name,
	}
	return v
}

func (v *Var) Name() string {
	return v.name
}

func (v *Var) String() string {
	return "$" + v.name
}

func (v *Var) GenVarExp(ps ...interface{}) ExpGenerator {
	return &Exp{
		TokenAndParams: append([]interface{}{v}, getTokenAndParams(ps)...),
	}
}

func (v *Var) WithFunction(f token.T, _ bool) *Exp {
	return &Exp{
		TokenAndParams: []interface{}{f, token.LPAREN, v, token.RPAREN},
	}
}

func (v *Var) Group(ps ...interface{}) *ExpGroup {
	return newExpGroup(nil, append([]interface{}{v}, ps...))
}

type Exp struct {
	from                *Exp
	TokenAndParams      []interface{}
	FirstTokenAndParams []interface{}
}

func (e *Exp) From() *Exp {
	return e.from
}

func (e *Exp) GenVarExp(ps ...interface{}) ExpGenerator {
	return &Exp{
		from:           e,
		TokenAndParams: getTokenAndParams(ps),
	}
}

func (e *Exp) WithBracket() {
	e.TokenAndParams = append([]interface{}{token.LPAREN}, e.TokenAndParams...)
	e.TokenAndParams = append(e.TokenAndParams, token.RPAREN)
}

func (e *Exp) WithFunction(f token.T, addBracket bool) *Exp {
	ds := []interface{}{f}
	if addBracket {
		ds = append(ds, token.LPAREN)
	}
	ds = append(ds, e.TokenAndParams...)
	if addBracket {
		ds = append(ds, token.RPAREN)
	}

	e.TokenAndParams = ds
	return e
}

func (e *Exp) Group(ps ...interface{}) *ExpGroup {
	return newExpGroup(e, ps...)
}

type ExpGroup struct {
	From  *Exp
	Exps  []ExpGenerator
	Items []interface{}
}

func newExpGroup(from *Exp, ps ...interface{}) *ExpGroup {
	eg := &ExpGroup{
		From: from,
	}

	eg.GenVarExp(ps...)
	return eg
}

func (e *ExpGroup) GenVarExp(ps ...interface{}) ExpGenerator {
	for _, p := range ps {
		if exp, ok := p.(ExpGenerator); ok {
			e.Exps = append(e.Exps, exp)
		} else {
			e.Items = append(e.Items, p)
		}
	}

	return e
}

func (e *ExpGroup) Group(ps ...interface{}) *ExpGroup {
	eg, _ := e.GenVarExp(ps...).(*ExpGroup)
	return eg
}

func (e *ExpGroup) WithFunction(t token.T, _ bool) *Exp {
	exp := &Exp{from: e.From, FirstTokenAndParams: []interface{}{t, token.LPAREN}}
	for i, ex := range e.Exps {
		if e1, ok := ex.(*Exp); ok {
			exp.TokenAndParams = append(exp.TokenAndParams, e1.TokenAndParams...)
		} else {
			exp.TokenAndParams = append(exp.TokenAndParams, ex)
		}
		if i < len(e.Exps)-1 {
			exp.TokenAndParams = append(exp.TokenAndParams, token.COMMA)
		}
	}

	if len(e.Items) > 0 {
		exp.TokenAndParams = append(exp.TokenAndParams, token.COMMA, e.Items[0])
	}
	exp.TokenAndParams = append(exp.TokenAndParams, token.RPAREN)
	return exp
}

func getTokenAndParams(ps []interface{}) []interface{} {
	r := make([]interface{}, 0)
	for _, p := range ps {
		if exp, ok := p.(*Exp); ok {
			r = append(r, exp.TokenAndParams...)
		} else {
			r = append(r, p)
		}
	}

	return r
}
