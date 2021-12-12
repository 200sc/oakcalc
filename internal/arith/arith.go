package arith

import "fmt"

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

func Eval(tokens []Token) (int64, error) {
	return 0, fmt.Errorf("unimplemented")
}

// 	if len(tokens) == 0 {
// 		return 0, nil
// 	}
// 	tk := tokens[0]
// 	// special case: - as negative sign
// 	if tk.Op != nil && *tk.Op == OpMinus {
// 		// TODO: tail recursion
// 		i, err := Eval(tokens[1:])
// 		if err != nil {
// 			return 0, err
// 		}
// 		return i * -1, nil
// 	}
// 	if tk.Op != nil {
// 		return 0, fmt.Errorf("expression began with binary operator")
// 	}
// 	return evalWithLHS(*tk.Number, tokens[1:])
// }

// func evalWithLHS(lhs int64, tokens []Token) (int64, error) {
// 	if len(tokens) == 0 {
// 		return lhs, nil
// 	}
// 	tk := tokens[0]
// 	i := 1
// 	for tk.Number != nil && i < len(tokens); i++ {
// 		lhs *= 10
// 		lhs += *tk.Number
// 		tk = tokens[i]
// 	}
// 	if i == len(tokens) {
// 		return lhs, nil
// 	}
// 	i++
// 	tk = tokens[i]
// 	switch
// }
