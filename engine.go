package golan

import (
	"fmt"
)

type Engine struct {
	env map[string]Value
}

func NewEngine() *Engine {
	return &Engine{
		env: map[string]Value{
			"print": NativeFunction(func(e *Engine, args []Value) (Value, error) {
				for _, v := range args {
					fmt.Println(v)
				}
				return Undefined{}, nil
			}),
		},
	}
}

func (e *Engine) Execute(tree Node) (Value, error) {
	return e.execNode(tree)
}

func (e *Engine) execNode(node Node) (Value, error) {
	switch n := node.(type) {
	case *Block:
		var r Value
		for _, c := range n.statements {
			v, err := e.execNode(c)
			if err != nil {
				return nil, err
			}
			r = v
		}
		return r, nil
	case *While:
		var r Value
		for {
			v, err := e.execNode(n.Condition)
			if err != nil {
				return nil, err
			}
			if !ValueTest(v) {
				break
			}
			v, err = e.execNode(n.Body)
			if err != nil {
				return nil, err
			}
			r = v
		}
		return r, nil
	case *If:
		v, err := e.execNode(n.Test)
		if err != nil {
			return nil, err
		}
		if ValueTest(v) {
			return e.execNode(n.Then)
		}
		if n.Alt == nil {
			return Undefined{}, nil
		}
		return e.execNode(n.Alt)
	case *Assign:
		v, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		e.env[n.Destination.(*Identifier).Name] = v
		return v, nil
	case *Equal:
		left, err := e.execNode(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.execNode(n.Right)
		if err != nil {
			return nil, err
		}
		return Boolean(left == right), nil
	case *NotEqual:
		left, err := e.execNode(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.execNode(n.Right)
		if err != nil {
			return nil, err
		}
		return Boolean(left != right), nil
	case *GreaterThanEqual:
		return e.execCompareOperation(n.Left, n.Right, func(x Integer, y Integer) bool { return x >= y })
	case *LessThanEqual:
		return e.execCompareOperation(n.Left, n.Right, func(x Integer, y Integer) bool { return x <= y })
	case *GreaterThan:
		return e.execCompareOperation(n.Left, n.Right, func(x Integer, y Integer) bool { return x > y })
	case *LessThan:
		return e.execCompareOperation(n.Left, n.Right, func(x Integer, y Integer) bool { return x < y })
	case *Addition:
		left, err := e.execNode(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.execNode(n.Right)
		if err != nil {
			return nil, err
		}
		return e.tryArithmeticOperation(left, n.Left.Position(), right, n.Right.Position(), func(x Integer, y Integer) Integer { return x + y })
	case *Subtraction:
		return e.execArithmeticOperation(n.Left, n.Right, func(x Integer, y Integer) Integer { return x - y })
	case *Multiplication:
		return e.execArithmeticOperation(n.Left, n.Right, func(x Integer, y Integer) Integer { return x * y })
	case *Division:
		left, err := e.execNode(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.execNode(n.Right)
		if err != nil {
			return nil, err
		}
		if IsInteger(left) && IsIntZero(right) {
			return nil, fmt.Errorf("%s: divided by zero", n.Position())
		}
		return e.tryArithmeticOperation(left, n.Left.Position(), right, n.Right.Position(), func(x Integer, y Integer) Integer { return x / y })
	case *Modulo:
		left, err := e.execNode(n.Left)
		if err != nil {
			return nil, err
		}
		right, err := e.execNode(n.Right)
		if err != nil {
			return nil, err
		}
		if IsInteger(left) && IsIntZero(right) {
			return nil, fmt.Errorf("%s: divided by zero", n.Position())
		}
		return e.tryArithmeticOperation(left, n.Left.Position(), right, n.Right.Position(), func(x Integer, y Integer) Integer { return x % y })
	case *Plus:
		val, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case Integer:
			return v, nil
		}
		return nil, fmt.Errorf("%s: cannot unary-plus operation for %v(%T)", n.Position(), val, val)
	case *Minus:
		val, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case Integer:
			return -v, nil
		}
		return nil, fmt.Errorf("%s: cannot unary-minus operation for %v(%T)", n.Position(), val, val)
	case *Not:
		v, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		return Boolean(!ValueTest(v)), nil
	case *IntLiteral:
		return Integer(n.Value), nil
	case *Identifier:
		v, ok := e.env[n.Name]
		if !ok {
			return nil, fmt.Errorf("%s: undefined variable - %s", n.position, n.Name)
		}
		return v, nil
	case *Apply:
		return e.execApply(n)
	}
	panic("must not happen")
}

func (e *Engine) execCompareOperation(left Node, right Node, intOp func(Integer, Integer) bool) (Value, error) {
	l, err := e.execNode(left)
	if err != nil {
		return nil, err
	}
	r, err := e.execNode(right)
	if err != nil {
		return nil, err
	}
	switch x := l.(type) {
	case Integer:
		switch y := r.(type) {
		case Integer:
			return Boolean(intOp(x, y)), nil
		default:
			return nil, fmt.Errorf("%s: incomparable type with Integer - %v(%T)", right.Position(), r, r)
		}
	}
	return nil, fmt.Errorf("%s: incomparable type - %v(%T)", left.Position(), l, l)
}

func (e *Engine) execArithmeticOperation(left Node, right Node, intOp func(Integer, Integer) Integer) (Value, error) {
	l, err := e.execNode(left)
	if err != nil {
		return nil, err
	}
	r, err := e.execNode(right)
	if err != nil {
		return nil, err
	}
	return e.tryArithmeticOperation(l, left.Position(), r, right.Position(), intOp)
}

func (e *Engine) tryArithmeticOperation(left Value, lp *Position, right Value, rp *Position, intOp func(Integer, Integer) Integer) (Value, error) {
	switch l := left.(type) {
	case Integer:
		switch r := right.(type) {
		case Integer:
			return intOp(l, r), nil
		default:
			return nil, fmt.Errorf("%s: not a Integer - %v(%T)", rp, right, right)
		}
	}
	return nil, fmt.Errorf("%s: not an addable type - %v(%T)", lp, left, left)
}

func (e *Engine) execApply(a *Apply) (Value, error) {
	v, err := e.execNode(a.function)
	if err != nil {
		return nil, err
	}
	f, ok := v.(NativeFunction)
	if !ok {
		return nil, fmt.Errorf("%s: not a function - %v(%T)", a.function.Position(), v, v)
	}
	args := []Value{}
	for _, x := range a.arguments {
		v, err := e.execNode(x)
		if err != nil {
			return nil, err
		}
		args = append(args, v)
	}
	return f(e, args)
}
