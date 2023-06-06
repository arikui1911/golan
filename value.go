package golan

type Value interface{}

func ValueTest(v Value) bool {
	b, ok := v.(Boolean)
	if !ok {
		return true
	}
	return bool(b)
}

type Undefined struct{}

func (Undefined) String() string { return "#<undefined>" }

func IsUndefined(v Value) bool {
	_, ok := v.(Undefined)
	return ok
}

type Boolean bool

type Integer int64

func IsInteger(v Value) bool {
	_, ok := v.(Integer)
	return ok
}

func IsIntZero(v Value) bool {
	i, ok := v.(Integer)
	return ok && i == 0
}

type Float float64

func IsFloat(v Value) bool {
	_, ok := v.(Float)
	return ok
}

type NativeFunction func(*Engine, []Value) (Value, error)
