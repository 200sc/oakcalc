package arith

import (
	"fmt"
	"io"
	"math"
	"strconv"
)

type Token struct {
	// One of:
	Number *int64
	Op     *Op
}

func (t Token) Copy() Token {
	t2 := Token{}
	if t.Number != nil {
		t2.Number = new(int64)
		*t2.Number = *t.Number
	}
	if t.Op != nil {
		t2.Op = new(Op)
		*t2.Op = *t.Op
	}
	return t2
}

type Op string

const (
	OpPlus       Op = "+"
	OpMinus      Op = "-"
	OpDivide     Op = "/"
	OpMultiply   Op = "*"
	OpEquals     Op = "="
	OpBackspace  Op = "<-"
	OpOpenParen  Op = "("
	OpCloseParen Op = ")"
	OpSquareRoot Op = "√"
)

var (
	unaryOps = map[Op]struct{}{
		OpSquareRoot: {},
		OpMinus:      {},
	}
	binaryOps = map[Op]struct{}{
		OpPlus:     {},
		OpDivide:   {},
		OpMultiply: {},
		OpMinus:    {},
	}
)

func (o Op) IsBinary() bool {
	_, ok := binaryOps[o]
	return ok
}

func (o Op) IsUnary() bool {
	_, ok := unaryOps[o]
	return ok
}

type Node interface {
	isNode()
}

type NumberNode int64

func (n NumberNode) isNode() {}

type BinaryOpNode struct {
	LHS, RHS Node
	Op
}

func (n BinaryOpNode) isNode() {}

type UnaryOpNode struct {
	Inner Node
	Op
}

func (n UnaryOpNode) isNode() {}

type ParenWrappedNode struct {
	Inner Node
}

func (n ParenWrappedNode) isNode() {}

func Eval(n Node) int64 {
	// assumes a well formed tree
	switch v := n.(type) {
	case NumberNode:
		return int64(v)
	case BinaryOpNode:
		lhs := Eval(v.LHS)
		rhs := Eval(v.RHS)
		switch v.Op {
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
	case UnaryOpNode:
		inner := Eval(v.Inner)
		switch v.Op {
		case OpSquareRoot:
			//  TODO: math.BigFloat
			return int64(math.Sqrt(float64(inner)))
		case OpMinus:
			return inner * -1
		default:
			return 0
		}
	case ParenWrappedNode:
		return Eval(v.Inner)
	default:
		panic("invalid node")
	}
}

func Pretty(n Node) string {
	switch v := n.(type) {
	case NumberNode:
		return strconv.FormatInt(int64(v), 10)
	case BinaryOpNode:
		lhs := Pretty(v.LHS)
		rhs := Pretty(v.RHS)
		return lhs + " " + string(v.Op) + " " + rhs
	case UnaryOpNode:
		return string(v.Op) + Pretty(v.Inner)
	case ParenWrappedNode:
		return "(" + Pretty(v.Inner) + ")"
	default:
		panic(fmt.Sprintf("invalid node: %T", n))
	}
}

func Parse(tokens []Token) (tree Node, err error) {
	var lhs Node
	if len(tokens) == 0 {
		return nil, io.EOF
	}
	i := 0
	thisToken := tokens[i]
	if len(tokens) == 1 {
		if thisToken.Number == nil {
			return nil, fmt.Errorf("expected number as final token")
		}
		lhs = NumberNode(*thisToken.Number)
		return lhs, nil
	}
	switch {
	case thisToken.Number != nil: // eq = number
		lhs = NumberNode(*thisToken.Number)
	case thisToken.Op != nil && thisToken.Op.IsUnary():
		i++
		nextToken := tokens[i]
		// if we call parse, how do we stop parse from consuming the entire rest of the tree?
		var inner Node
		switch {
		// TODO -()
		case nextToken.Op != nil && *nextToken.Op == OpOpenParen:
			i++
			inner, i, err = parseParenExpr(i, tokens)
			if err != nil {
				return nil, err
			}
		case nextToken.Number != nil:
			inner = NumberNode(*nextToken.Number)
		default:
			return nil, fmt.Errorf("expected number after -")
		}
		lhs = UnaryOpNode{
			Inner: inner,
			Op:    *thisToken.Op,
		}
	case thisToken.Op != nil && *thisToken.Op == OpOpenParen:
		i++
		lhs, i, err = parseParenExpr(i, tokens)
		if err != nil {
			return lhs, err
		}
	default:
		return lhs, fmt.Errorf("bad token for starting production")
	}
	i++
	if i >= len(tokens) {
		return lhs, nil
	}
	nextToken := tokens[i]
	// based on our productions, next token has to be a binop.
	if nextToken.Op == nil {
		return tree, fmt.Errorf("eq middle token %v was not an op", nextToken)
	}
	newTree := BinaryOpNode{
		LHS: lhs,
		Op:  *nextToken.Op,
		RHS: nil,
	}
	i++
	right, err := Parse(tokens[i:])
	if err != nil {
		return newTree, err
	}
	newTree.RHS = right
	return newTree, nil
}

// parens are the worst
// we're doing a bad job here
// read until we get a close paren
func parseParenExpr(j int, tokens []Token) (node Node, newI int, err error) {
	startJ := j
	closeNeeded := 1
	for j < len(tokens) {
		if tokens[j].Op != nil && *tokens[j].Op == OpOpenParen {
			closeNeeded++
		}
		if tokens[j].Op != nil && *tokens[j].Op == OpCloseParen {
			closeNeeded--
			if closeNeeded == 0 {
				break
			}
		}
		j++
	}
	inner, err := Parse(tokens[startJ:j])
	if err != nil {
		return nil, 0, err
	}
	return ParenWrappedNode{
		Inner: inner,
	}, j, nil
}

func ParseString(s string) (Node, error) {
	tks := []Token{}
	for _, c := range s {
		if c == ' ' {
			continue
		}
		switch c {
		case '0':
			tks = append(tks, iTk(0))
		case '1':
			tks = append(tks, iTk(1))
		case '2':
			tks = append(tks, iTk(2))
		case '3':
			tks = append(tks, iTk(3))
		case '4':
			tks = append(tks, iTk(4))
		case '5':
			tks = append(tks, iTk(5))
		case '6':
			tks = append(tks, iTk(6))
		case '7':
			tks = append(tks, iTk(7))
		case '8':
			tks = append(tks, iTk(8))
		case '9':
			tks = append(tks, iTk(9))
		case '(':
			tks = append(tks, oTk(OpOpenParen))
		case ')':
			tks = append(tks, oTk(OpCloseParen))
		case '-':
			tks = append(tks, oTk(OpMinus))
		case '+':
			tks = append(tks, oTk(OpPlus))
		case '*':
			tks = append(tks, oTk(OpMultiply))
		case '/':
			tks = append(tks, oTk(OpDivide))
		case '√':
			tks = append(tks, oTk(OpSquareRoot))
		}
		if len(tks) > 1 {
			if tks[len(tks)-1].Number != nil &&
				tks[len(tks)-2].Number != nil {
				dec := *tks[len(tks)-2].Number
				dig := *tks[len(tks)-1].Number
				tks[len(tks)-2].Number = i64p(dec*10 + dig)
				tks = tks[:len(tks)-1]
			}
		}
	}
	return Parse(tks)
}

// productions
// eq =
//   eq binop eq
//   ( eq )
//   unop eq
//   numeral
// binop = - | + | * | /
// unop = - | √
//
// start = eq

func iTk(i int64) Token {
	return Token{
		Number: i64p(i),
	}
}
func oTk(o Op) Token {
	return Token{
		Op: opP(o),
	}
}

func i64p(i int64) *int64 {
	return &i
}

func opP(o Op) *Op {
	return &o
}
