package microgorod_test

import (
	"fmt"
	"reflect"
	"testing"

	mg "github.com/mkmik/microgorod"
)

func TestForward(t *testing.T) {
	a := mg.Value(10.0)
	b := mg.Value(2.0)

	testCases := []struct {
		forward func() *mg.Expr
		want    float64
	}{
		{
			func() *mg.Expr { return mg.Add(a, b) },
			12.0,
		},
		{
			func() *mg.Expr { return mg.Mul(a, b) },
			20.0,
		},
		{
			func() *mg.Expr { return mg.Pow(a, b.Data) },
			100.0,
		},
		{
			func() *mg.Expr { return mg.Sub(a, b) },
			8.0,
		},
		{
			func() *mg.Expr { return mg.Div(a, b) },
			5.0,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			got := tc.forward()
			if got, want := got.Data, tc.want; got != want {
				t.Fatalf("got: %v, want: %v", got, want)
			}
		})
	}
}

func TestBackward(t *testing.T) {
	testCases := []struct {
		params  []*mg.Expr
		forward func(...*mg.Expr) *mg.Expr
		want    []float64
	}{
		{
			[]*mg.Expr{mg.Value(10.0), mg.Value(20.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Add(ps[0], ps[1])
			},
			[]float64{1.0, 1.0},
		},
		{
			[]*mg.Expr{mg.Value(10.0), mg.Value(20.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Sub(ps[0], ps[1])
			},
			[]float64{1.0, -1.0},
		},
		{
			[]*mg.Expr{mg.Value(10.0), mg.Value(20.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Mul(ps[0], ps[1])
			},
			[]float64{20.0, 10.0},
		},
		{
			[]*mg.Expr{mg.Value(10.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Pow(ps[0], 2)
			},
			[]float64{20.0},
		},
		{
			[]*mg.Expr{mg.Value(10.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Pow(ps[0], 3)
			},
			[]float64{300.0},
		},
		{
			[]*mg.Expr{mg.Value(10.0), mg.Value(2.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Div(ps[0], ps[1])
			},
			[]float64{0.5, -2.5},
		},
		{
			[]*mg.Expr{mg.Value(10.0), mg.Value(20.0), mg.Value(2.0)},
			func(ps ...*mg.Expr) *mg.Expr {
				return mg.Add(ps[0], mg.Mul(ps[1], ps[2]))
			},
			[]float64{1.0, 2.0, 20.0},
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			res := tc.forward(tc.params...)
			res.Backward()

			var grads []float64
			for _, p := range tc.params {
				grads = append(grads, p.Grad)
			}

			if got, want := grads, tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("got: %v, want: %v", got, want)
			}
		})
	}
}
