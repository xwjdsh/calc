package token

type T string

func (t T) String() string {
	return string(t)
}

func (t T) Token() T {
	return t
}

func (t T) Pointer() *T {
	return &t
}

const (
	// general type
	ADD   T = "+"
	SUB   T = "-"
	MUL   T = "*"
	QUO   T = "/"
	REM   T = "%"
	COMMA T = ","

	// bracket type
	LPAREN T = "("
	RPAREN T = ")"

	// function type
	SIN T = "sin"
	COS T = "cos"
	TAN T = "tan"
	ABS T = "abs"
	OPP T = "opp" // opposite number
	REC T = "rec" // reciprocal
	SUM T = "sum"
	MAX T = "max"
	MIN T = "min"
	POW T = "pow"
)
