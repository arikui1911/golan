package golan

import (
	"fmt"
	"io"
)

type Position struct {
	FirstLineno int
	FirstColumn int
	LastLineno  int
	LastColumn  int
}

func (p *Position) String() string {
	return fmt.Sprintf(
		"(%d:%d,%d:%d)",
		p.FirstLineno, p.FirstColumn,
		p.LastLineno, p.LastColumn,
	)
}

func DumpTree(tree Node, output io.Writer) {
	tree.dump(output, 0)
}

func indent(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		fmt.Fprint(w, "  ")
	}
}

type Node interface {
	Position() *Position
	dump(output io.Writer, indentLevel int)
}

type Assign struct {
	position    *Position
	Destination Node
	Expression  Node
}

func (a *Assign) Position() *Position { return a.position }

func (a *Assign) dump(w io.Writer, n int) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v\n", a, a.position)
	a.Destination.dump(w, n+1)
	a.Expression.dump(w, n+1)
}

func dumpBinary(w io.Writer, n int, node Node, l Node, r Node) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v\n", node, node.Position())
	l.dump(w, n+1)
	r.dump(w, n+1)
}

type Addition struct {
	position *Position
	Left     Node
	Right    Node
}

func (a *Addition) Position() *Position { return a.position }

func (a *Addition) dump(w io.Writer, n int) {
	dumpBinary(w, n, a, a.Left, a.Right)
}

type Subtraction struct {
	position *Position
	Left     Node
	Right    Node
}

func (s *Subtraction) Position() *Position { return s.position }

func (s *Subtraction) dump(w io.Writer, n int) {
	dumpBinary(w, n, s, s.Left, s.Right)
}

type Multiplication struct {
	position *Position
	Left     Node
	Right    Node
}

func (m *Multiplication) Position() *Position { return m.position }

func (m *Multiplication) dump(w io.Writer, n int) {
	dumpBinary(w, n, m, m.Left, m.Right)
}

type Division struct {
	position *Position
	Left     Node
	Right    Node
}

func (d *Division) Position() *Position { return d.position }

func (d *Division) dump(w io.Writer, n int) {
	dumpBinary(w, n, d, d.Left, d.Right)
}

type Modulo struct {
	position *Position
	Left     Node
	Right    Node
}

func (m *Modulo) Position() *Position { return m.position }

func (m *Modulo) dump(w io.Writer, n int) {
	dumpBinary(w, n, m, m.Left, m.Right)
}

type IntLiteral struct {
	position *Position
	Value    int64
}

func (i *IntLiteral) Position() *Position { return i.position }

func (i *IntLiteral) dump(w io.Writer, n int) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v: %v\n", i, i.position, i.Value)
}

type Identifier struct {
	position *Position
	Name     string
}

func (i *Identifier) Position() *Position { return i.position }

func (i *Identifier) dump(w io.Writer, n int) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v: %v\n", i, i.position, i.Name)
}
