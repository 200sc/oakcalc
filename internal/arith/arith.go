package arith

import (
	"fmt"
	"strconv"
	"strings"
)

type Token struct {
	// One of:
	Number *int64
	Op     *Op
}

type Op string

const (
	OpPlus      Op = "+"
	OpMinus     Op = "-"
	OpDivide    Op = "/"
	OpMultiply  Op = "*"
	OpEquals    Op = "="
	OpBackspace Op = "<-"
)

type Tree struct {
	Token
	Left  *Tree
	Right *Tree
}

func (t Tree) Eval() int64 {
	// assumes a well formed tree

	if t.Left != nil && t.Right != nil {
		lhs := t.Left.Eval()
		rhs := t.Right.Eval()
		switch *t.Op {
		case OpDivide:
			// todo: divide by zero detection
			return lhs / rhs
		case OpMultiply:
			return lhs * rhs
		case OpMinus:
			return lhs - rhs
		case OpPlus:
			return lhs + rhs
		}
		return 0
	}
	return *t.Number
}

func (t Tree) Pretty() string {
	sb := &strings.Builder{}
	t.pretty(sb)
	return sb.String()
}

func (t Tree) pretty(sb *strings.Builder) {
	if t.Left != nil {
		t.Left.pretty(sb)
		sb.WriteString(" ")
	}
	if t.Token.Number != nil {
		sb.WriteString(strconv.FormatInt(*t.Token.Number, 10))
	}
	if t.Token.Op != nil {
		sb.WriteString(string(*t.Token.Op))
	}
	if t.Right != nil {
		sb.WriteString(" ")
		t.Right.pretty(sb)
	}
}

func Parse(tokens []Token) (*Tree, error) {
	tree := &Tree{}
	if len(tokens) == 0 {
		return tree, nil
	}
	tree.Token = tokens[0]
	i := 1
	for {
		if i >= len(tokens) {
			return tree, nil
		}
		tk := tokens[i]
		if tk.Op != nil {
			// Assuming binary op
			if i+1 >= len(tokens) {
				return nil, fmt.Errorf("unpaired binary operator")
			}
			left := tree
			right, err := Parse(tokens[i+1:])
			if err != nil {
				return nil, err
			}
			tree = &Tree{
				Token: tk,
				Left:  left,
				Right: right,
			}
		}
		i++
	}
}
