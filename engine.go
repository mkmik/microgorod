package microgorod

import (
	"fmt"
	"math"
	"slices"
)

type Expr struct {
	Op   string
	Data float64
	Grad float64

	Children []*Expr
	back     func(*Expr)
}

func New(op string, v float64, back func(*Expr), c ...*Expr) *Expr {
	return &Expr{
		Op:       op,
		Data:     v,
		Children: c,
		back:     back,
	}
}

func Value(v float64) *Expr {
	return New("n", v, func(*Expr) {})
}

// fundamental ops

func Add(a, b *Expr) *Expr {
	back := func(self *Expr) {
		self.chain(a, 1)
		self.chain(b, 1)
	}
	return New("+", a.Data+b.Data, back, a, b)
}

func Mul(a, b *Expr) *Expr {
	back := func(self *Expr) {
		self.chain(a, b.Data)
		self.chain(b, a.Data)
	}
	return New("*", a.Data*b.Data, back, a, b)
}

func Pow(a *Expr, e float64) *Expr {
	back := func(self *Expr) {
		self.chain(a, e*math.Pow(a.Data, e-1))
	}
	return New("^", math.Pow(a.Data, e), back, a)
}

// derived ops

func Sub(a, b *Expr) *Expr {
	return Add(a, Mul(b, Value(-1.0)))
}

func Div(a, b *Expr) *Expr {
	return Mul(a, Pow(b, -1.0))
}

// backward pass

func (e *Expr) Backward() {
	e.Grad = 1
	for _, e := range slices.Backward(e.topo()) {
		e.back(e)
	}
}

// helpers

// chain rule
func (e *Expr) chain(a *Expr, v float64) {
	a.Grad += v * e.Grad
}

// topological sort
func (e *Expr) topo() []*Expr {
	var (
		res  []*Expr
		seen = map[*Expr]bool{}
		rec  func(*Expr)
	)
	rec = func(e *Expr) {
		if !seen[e] {
			seen[e] = true
			for _, c := range e.Children {
				rec(c)
			}
			res = append(res, e)
		}
	}
	rec(e)
	return res
}

func (e *Expr) GoString() string {
	return fmt.Sprintf("%s(%v, ∆%v, %d children)", e.Op, e.Data, e.Grad, len(e.Children))
}
