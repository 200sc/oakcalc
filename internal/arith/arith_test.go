package arith

import (
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		in  string
		out int64
	}
	tcs := []testCase{
		{
			in:  "0+0",
			out: 0,
		}, {
			in:  "(5*3)-1",
			out: 14,
		}, {
			in:  "1 + (3 + (2 + (9)))",
			out: 15,
		}, {
			in:  "-50",
			out: -50,
		},
	}
	for i, tc := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tree, _, err := ParseString(tc.in)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}
			res := Eval(tree)
			if res != tc.out {
				t.Fatalf("out mismatch: expected %v vs %v", tc.out, res)
			}
		})
	}
}
