package microgorod

import (
	"fmt"
	"testing"
)

func TestTopo(t *testing.T) {
	testCases := []struct {
		tree *Expr
		want []*Expr
	}{
		{
			Add(Value(1), Value(2)),
			[]*Expr{{Op: "n", Data: 1}, {Op: "n", Data: 2}, {Op: "+", Data: 3}},
		},
		{
			Add(Value(1), Mul(Value(2), Value(3))),
			[]*Expr{{Op: "n", Data: 1}, {Op: "n", Data: 2}, {Op: "n", Data: 3}, {Op: "*", Data: 6}, {Op: "+", Data: 7}},
		},
	}

	equal := func(a, b []*Expr) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range a {
			if a[i].Data != b[i].Data {
				return false
			}
			if a[i].Op != b[i].Op {
				return false
			}
		}
		return true
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if got, want := tc.tree.topo(), tc.want; !equal(got, want) {
				t.Fatalf("got: %#v, want: %#v", got, want)
			}
		})
	}
}
