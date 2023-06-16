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
		return e.execCmp(n.Left, n.Right, n.Position(), CMP_EQ)
	case *NotEqual:
		r, err := e.execCmp(n.Left, n.Right, n.Position(), CMP_EQ)
		if err != nil {
			return nil, err
		}
		return !(r.(Boolean)), nil
	case *GreaterThanEqual:
		return e.execCmp(n.Left, n.Right, n.Position(), CMP_GREATER, CMP_EQ)
	case *LessThanEqual:
		return e.execCmp(n.Left, n.Right, n.Position(), CMP_LESS, CMP_EQ)
	case *GreaterThan:
		return e.execCmp(n.Left, n.Right, n.Position(), CMP_GREATER)
	case *LessThan:
		return e.execCmp(n.Left, n.Right, n.Position(), CMP_LESS)
	case *Addition:
		return e.execBinArith(n.Left, n.Right, n.Position(), AddValues)
	case *Subtraction:
		return e.execBinArith(n.Left, n.Right, n.Position(), SubtractValues)
	case *Multiplication:
		return e.execBinArith(n.Left, n.Right, n.Position(), MultiplyValues)
	case *Division:
		return e.execBinArith(n.Left, n.Right, n.Position(), DivideValues)
	case *Modulo:
		return e.execBinArith(n.Left, n.Right, n.Position(), ModuloValues)
	case *Plus:
		val, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		v, ok := val.(SignableValue)
		if !ok {
			return nil, fmt.Errorf("%s: invalid plus sign with %v(%T)", n.Position(), val, val)
		}
		return v.OpPlus()
	case *Minus:
		val, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		v, ok := val.(SignableValue)
		if !ok {
			return nil, fmt.Errorf("%s: invalid minus sign with %v(%T)", n.Position(), val, val)
		}
		return v.OpMinus()
	case *Not:
		v, err := e.execNode(n.Expression)
		if err != nil {
			return nil, err
		}
		return Boolean(!ValueTest(v)), nil
	case *BooleanLiteral:
		return Boolean(n.Value), nil
	case *IntLiteral:
		return Integer(n.Value), nil
	case *FloatLiteral:
		return Float(n.Value), nil
	case *StringLiteral:
		return String(n.Value), nil
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

func (e *Engine) execCmp(left Node, right Node, p *Position, wants ...CompareResult) (Value, error) {
	l, err := e.execNode(left)
	if err != nil {
		return nil, err
	}
	r, err := e.execNode(right)
	if err != nil {
		return nil, err
	}
	result, err := CompareValues(l, r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	for _, want := range wants {
		if result == want {
			return Boolean(true), nil
		}
	}
	return Boolean(false), nil
}

func (e *Engine) execBinArith(left Node, right Node, p *Position, op func(Value, Value) (Value, error)) (Value, error) {
	l, err := e.execNode(left)
	if err != nil {
		return nil, err
	}
	r, err := e.execNode(right)
	if err != nil {
		return nil, err
	}
	result, err := op(l, r)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", p, err)
	}
	return result, nil
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
