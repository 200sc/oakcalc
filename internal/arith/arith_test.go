package arith

import (
	"fmt"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	type testCase struct {
		tokens []Token
		tree   *Tree
	}
	tcs := []testCase{
		{
			tokens: []Token{
				ntk(0),
				otk(OpPlus),
				ntk(0),
			},
			tree: &Tree{
				Left:  &Tree{Token: ntk(0)},
				Token: otk(OpPlus),
				Right: &Tree{Token: ntk(0)},
			},
		},
	}
	for i, tc := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tree, _ := Parse(tc.tokens)
			if err := TreesEqual(tree, tc.tree); err != nil {
				t.Errorf("mismatched tree: %v", err)
				fmt.Println(tree.Pretty())
				fmt.Println(tc.tree.Pretty())
			}
		})
	}
}

func TreesEqual(a, b *Tree) error {
	if (a.Number == nil) != (b.Number == nil) {
		return fmt.Errorf("number nil mismatch")
	}
	if a.Number != nil {
		if *a.Number != *b.Number {
			return fmt.Errorf("number mismatch (%v vs %v)", *a.Number, *b.Number)
		}
	}
	if (a.Op == nil) != (b.Op == nil) {
		return fmt.Errorf("op nil mismatch")
	}
	if a.Op != nil {
		if *a.Op != *b.Op {
			return fmt.Errorf("op mismatch (%v vs %v)", *a.Op, *b.Op)
		}
	}
	if (a.Left == nil) != (b.Left == nil) {
		return fmt.Errorf("left nil mismatch")
	}
	if a.Left != nil {
		if err := TreesEqual(a.Left, b.Left); err != nil {
			return fmt.Errorf("left tree mismatch: %v", err)
		}
	}
	if (a.Right == nil) != (b.Right == nil) {
		return fmt.Errorf("right nil mismatch")
	}
	if a.Right != nil {
		if err := TreesEqual(a.Right, b.Right); err != nil {
			return fmt.Errorf("right tree mismatch: %v", err)
		}
	}
	return nil
}

func i64p(i int64) *int64 {
	return &i
}

func opP(o Op) *Op {
	return &o
}

func ntk(i int64) Token {
	return Token{Number: i64p(i)}
}

func otk(o Op) Token {
	return Token{Op: opP(o)}
}
