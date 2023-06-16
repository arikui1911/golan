package golan

import (
	"errors"
	"fmt"
)

type Value interface{}

func ValueTest(v Value) bool {
	b, ok := v.(Boolean)
	if !ok {
		return true
	}
	return bool(b)
}

type CompareResult int

const (
	CMP_EQ CompareResult = iota
	CMP_LESS
	CMP_GREATER
	CMP_INVALID
)

type ComparableValue interface {
	Value
	OpCmp(Value) CompareResult
}

func CompareValues(x Value, y Value) (CompareResult, error) {
	if a, ok := x.(ComparableValue); ok {
		r := a.OpCmp(y)
		if r != CMP_INVALID {
			return r, nil
		}
	}
	if x == y {
		return CMP_EQ, nil
	}
	return CMP_INVALID, fmt.Errorf("not a comparable value - %v(%T)", x, x)
}

type AddableValue interface {
	Value
	OpAdd(Value) (Value, error)
}

func AddValues(x Value, y Value) (Value, error) {
	if a, ok := x.(AddableValue); ok {
		return a.OpAdd(y)
	}
	return nil, fmt.Errorf("not an addable value - %v(%T)", x, x)
}

type ArithmeticValue interface {
	AddableValue
	OpSub(Value) (Value, error)
	OpMul(Value) (Value, error)
	OpDiv(Value) (Value, error)
}

func SubtractValues(x Value, y Value) (Value, error) {
	if a, ok := x.(ArithmeticValue); ok {
		return a.OpSub(y)
	}
	return nil, fmt.Errorf("not a subtractable value - %v(%T)", x, x)
}

func MultiplyValues(x Value, y Value) (Value, error) {
	if a, ok := x.(ArithmeticValue); ok {
		return a.OpMul(y)
	}
	return nil, fmt.Errorf("not a multipliable value - %v(%T)", x, x)
}

func DivideValues(x Value, y Value) (Value, error) {
	if a, ok := x.(ArithmeticValue); ok {
		return a.OpDiv(y)
	}
	return nil, fmt.Errorf("not a dividable value - %v(%T)", x, x)
}

type IntegerValue interface {
	ArithmeticValue
	OpMod(Value) (Value, error)
}

func ModuloValues(x Value, y Value) (Value, error) {
	if a, ok := x.(IntegerValue); ok {
		return a.OpMod(y)
	}
	return nil, fmt.Errorf("not a modulo-operatable value - %v(%T)", x, x)
}

type SignableValue interface {
	Value
	OpPlus() (Value, error)
	OpMinus() (Value, error)
}

type Undefined struct{}

func (Undefined) String() string { return "#<undefined>" }

func IsUndefined(v Value) bool {
	_, ok := v.(Undefined)
	return ok
}

type Boolean bool

type Integer int64

func (x Integer) OpCmp(other Value) CompareResult {
	y, ok := other.(Integer)
	if !ok {
		return CMP_INVALID
	}
	if x < y {
		return CMP_LESS
	}
	if x > y {
		return CMP_GREATER
	}
	return CMP_EQ
}

func (x Integer) OpAdd(other Value) (Value, error) {
	y, ok := other.(Integer)
	if !ok {
		return nil, fmt.Errorf("not a Integer - %v(%T)", other, other)
	}
	return x + y, nil
}

func (x Integer) OpSub(other Value) (Value, error) {
	y, ok := other.(Integer)
	if !ok {
		return nil, fmt.Errorf("not a Integer - %v(%T)", other, other)
	}
	return x - y, nil
}

func (x Integer) OpMul(other Value) (Value, error) {
	y, ok := other.(Integer)
	if !ok {
		return nil, fmt.Errorf("not a Integer - %v(%T)", other, other)
	}
	return x * y, nil
}

func (x Integer) OpDiv(other Value) (Value, error) {
	y, ok := other.(Integer)
	if !ok {
		return nil, fmt.Errorf("not a Integer - %v(%T)", other, other)
	}
	if y == 0 {
		return nil, errors.New("divide by zero")
	}
	return x / y, nil
}

func (x Integer) OpMod(other Value) (Value, error) {
	y, ok := other.(Integer)
	if !ok {
		return nil, fmt.Errorf("not a Integer - %v(%T)", other, other)
	}
	if y == 0 {
		return nil, errors.New("divide by zero")
	}
	return x % y, nil
}

func (x Integer) OpPlus() (Value, error) { return x, nil }

func (x Integer) OpMinus() (Value, error) { return -x, nil }

type Float float64

func (x Float) OpAdd(other Value) (Value, error) {
	y, ok := other.(Float)
	if !ok {
		return nil, fmt.Errorf("not a Float - %v(%T)", other, other)
	}
	return x + y, nil
}

func (x Float) OpSub(other Value) (Value, error) {
	y, ok := other.(Float)
	if !ok {
		return nil, fmt.Errorf("not a Float - %v(%T)", other, other)
	}
	return x - y, nil
}

func (x Float) OpMul(other Value) (Value, error) {
	y, ok := other.(Float)
	if !ok {
		return nil, fmt.Errorf("not a Float - %v(%T)", other, other)
	}
	return x * y, nil
}

func (x Float) OpDiv(other Value) (Value, error) {
	y, ok := other.(Float)
	if !ok {
		return nil, fmt.Errorf("not a Float - %v(%T)", other, other)
	}
	if y == 0 {
		return nil, errors.New("divide by zero")
	}
	return x / y, nil
}

func (x Float) OpPlus() (Value, error) { return x, nil }

func (x Float) OpMinus() (Value, error) { return -x, nil }

type String string

func (x String) OpAdd(other Value) (Value, error) {
	y, ok := other.(String)
	if !ok {
		return nil, fmt.Errorf("not a String - %v(%T)", other, other)
	}
	return x + y, nil
}

type NativeFunction func(*Engine, []Value) (Value, error)
