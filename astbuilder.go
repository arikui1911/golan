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
