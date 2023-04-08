package golan

import (
	"errors"
	"strconv"
	"strings"
)

type ASTBuilder struct {
	buffer  string
	stack   []Node
	lastErr error
}

func (b *ASTBuilder) ASTBuilderInit(buffer string) {
	b.buffer = buffer
	b.push(&Block{[]Node{}})
}

func (b *ASTBuilder) Err() error { return b.lastErr }

type bailout struct{}

func (b *ASTBuilder) Raise(e error) {
	b.lastErr = e
	panic(bailout{})
}

func (b *ASTBuilder) Recover(e any) {
	if e == nil {
		return
	}
	if _, ok := e.(bailout); !ok {
		panic(e)
	}
}

func (b *ASTBuilder) Finish() Node {
	if len(b.stack) == 0 {
		return nil
	}
	return b.pop()
}

func (b *ASTBuilder) push(n Node) {
	b.stack = append(b.stack, n)
}

func (b *ASTBuilder) pop() Node {
	if len(b.stack) == 0 {
		b.Raise(errors.New("cannot pop from empty stack"))
	}
	r := b.stack[len(b.stack)-1]
	b.stack = b.stack[:len(b.stack)-1]
	return r
}

func (b *ASTBuilder) PushBlock() {
	b.push(&Block{[]Node{}})
}

func (b *ASTBuilder) PushWhile(beg int) {
	fl, fc := calcPosition(b.buffer, beg)
	b.push(&While{position: &Position{fl, fc, 0, 0}})
}

func (b *ASTBuilder) CompleteWhile() {
	body := b.pop()
	cond := b.pop()
	w := b.pop().(*While)
	current := b.pop().(*Block)
	w.position.LastLineno = body.Position().LastLineno
	w.position.LastColumn = body.Position().LastColumn
	w.Condition = cond
	w.Body = body
	current.Add(w)
	b.push(current)
}

func (b *ASTBuilder) PushExpressionStatement() {
	x := b.pop()
	block := b.pop().(*Block)
	block.Add(x)
	b.push(block)
}

func (b *ASTBuilder) PushAssign() {
	r := b.pop()
	l := b.pop()
	p := &Position{
		l.Position().FirstLineno, l.Position().FirstColumn,
		r.Position().LastLineno, r.Position().LastColumn,
	}
	b.push(&Assign{p, l, r})
}

func (b *ASTBuilder) PushBinOp(op string) {
	y := b.pop()
	x := b.pop()
	p := &Position{
		x.Position().FirstLineno, x.Position().FirstColumn,
		y.Position().LastLineno, y.Position().LastColumn,
	}
	switch op {
	case "==":
		b.push(&Equal{p, x, y})
	case "!=":
		b.push(&NotEqual{p, x, y})
	case ">=":
		b.push(&GreaterThanEqual{p, x, y})
	case "<=":
		b.push(&LessThanEqual{p, x, y})
	case ">":
		b.push(&GreaterThan{p, x, y})
	case "<":
		b.push(&LessThan{p, x, y})
	case "+":
		b.push(&Addition{p, x, y})
	case "-":
		b.push(&Subtraction{p, x, y})
	case "*":
		b.push(&Multiplication{p, x, y})
	case "/":
		b.push(&Division{p, x, y})
	case "%":
		b.push(&Modulo{p, x, y})
	default:
		panic("must not happen")
	}
}

func (b *ASTBuilder) PushUnaryOp(beg int, _ int, op string) {
	fl, fc := calcPosition(b.buffer, beg)
	switch op {
	case "+":
		b.push(&Plus{&Position{fl, fc, 0, 0}, nil})
	case "-":
		b.push(&Minus{&Position{fl, fc, 0, 0}, nil})
	case "!":
		b.push(&Not{&Position{fl, fc, 0, 0}, nil})
	default:
		panic("must not happen")
	}
}

func (b *ASTBuilder) CompleteUnary() {
	x := b.pop()
	u := b.pop()
	p := u.Position()
	p.LastLineno = x.Position().LastLineno
	p.LastColumn = x.Position().LastColumn
	u.(Unary).Complete(p, x)
	b.push(u)
}

func (b *ASTBuilder) PushIntLiteral(beg int, end int, src string) {
	fl, fc := calcPosition(b.buffer, beg)
	ll, lc := calcPosition(b.buffer, end-1)
	i64, _ := strconv.ParseInt(src, 10, 64)
	b.push(&IntLiteral{&Position{fl, fc, ll, lc}, i64})
}

func (b *ASTBuilder) PushIdentifier(beg int, end int, src string) {
	fl, fc := calcPosition(b.buffer, beg)
	ll, lc := calcPosition(b.buffer, end-1)
	b.push(&Identifier{&Position{fl, fc, ll, lc}, src})
}

func calcPosition(src string, pos int) (int, int) {
	a := strings.Split(src[:pos], "\n")
	lineno := len(a) - 1
	column := len(a[lineno])
	return lineno, column
}
