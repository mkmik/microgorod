package main

import (
	"fmt"

	mg "github.com/mkmik/microgorod"
)

func main() {
	a := mg.Value(10.0)
	b := mg.Value(2.0)
	f := mg.Value(-2.0)

	params := []*mg.Expr{a, b, f}

	forward := func() *mg.Expr {
		c := mg.Add(a, b)
		return mg.Div(c, f)
	}
	g := forward()

	fmt.Printf("%.4f\n", g.Data)

	g.Backward()
	fmt.Printf("dg/da = %.4f\n", a.Grad)
	fmt.Printf("dg/db = %.4f\n", b.Grad)
	fmt.Printf("dg/df = %.4f\n", f.Grad)

	for i := 0; i < 100; i++ {
		for _, p := range params {
			p.Data += p.Grad * 0.1
		}

		g = forward()
		for _, p := range params {
			p.Grad = 0
		}
		g.Backward()
		fmt.Printf("%.2f\n", g.Data)
	}
}
