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
	if tree == nil {
		return
	}
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

type Block struct {
	statements []Node
}

func (b *Block) Position() *Position {
	var p Position
	if len(b.statements) > 0 {
		p.FirstLineno = b.statements[0].Position().FirstLineno
		p.FirstColumn = b.statements[0].Position().FirstColumn
		p.LastLineno = b.statements[len(b.statements)-1].Position().LastLineno
		p.LastColumn = b.statements[len(b.statements)-1].Position().LastColumn
	}
	return &p
}

func (b *Block) dump(w io.Writer, n int) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v\n", b, b.Position())
	for _, s := range b.statements {
		s.dump(w, n+1)
	}
}

func (b *Block) Add(statement Node) {
	b.statements = append(b.statements, statement)
}

type While struct {
	position  *Position
	Condition Node
	Body      Node
}

func (w *While) Position() *Position { return w.position }

func (w *While) dump(o io.Writer, n int) {
	indent(o, n)
	fmt.Fprintf(o, "%T:%v\n", w, w.position)
	w.Condition.dump(o, n+1)
	w.Body.dump(o, n+1)
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

type Equal struct {
	position *Position
	Left     Node
	Right    Node
}

func (e *Equal) Position() *Position { return e.position }

func (e *Equal) dump(w io.Writer, n int) {
	dumpBinary(w, n, e, e.Left, e.Right)
}

type NotEqual struct {
	position *Position
	Left     Node
	Right    Node
}

func (N *NotEqual) Position() *Position { return N.position }

func (N *NotEqual) dump(w io.Writer, n int) {
	dumpBinary(w, n, N, N.Left, N.Right)
}

type GreaterThanEqual struct {
	position *Position
	Left     Node
	Right    Node
}

func (G *GreaterThanEqual) Position() *Position { return G.position }

func (G *GreaterThanEqual) dump(w io.Writer, n int) {
	dumpBinary(w, n, G, G.Left, G.Right)
}

type LessThanEqual struct {
	position *Position
	Left     Node
	Right    Node
}

func (L *LessThanEqual) Position() *Position { return L.position }

func (L *LessThanEqual) dump(w io.Writer, n int) {
	dumpBinary(w, n, L, L.Left, L.Right)
}

type GreaterThan struct {
	position *Position
	Left     Node
	Right    Node
}

func (G *GreaterThan) Position() *Position { return G.position }

func (G *GreaterThan) dump(w io.Writer, n int) {
	dumpBinary(w, n, G, G.Left, G.Right)
}

type LessThan struct {
	position *Position
	Left     Node
	Right    Node
}

func (L *LessThan) Position() *Position { return L.position }

func (L *LessThan) dump(w io.Writer, n int) {
	dumpBinary(w, n, L, L.Left, L.Right)
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

type Unary interface {
	Node
	Complete(*Position, Node)
}

type Plus struct {
	position   *Position
	Expression Node
}

func (p *Plus) Complete(pos *Position, x Node) {
	p.position = pos
	p.Expression = x
}

func (p *Plus) Position() *Position { return p.position }

func (p *Plus) dump(w io.Writer, n int) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v:\n", p, p.position)
	p.Expression.dump(w, n+1)
}

type Minus struct {
	position   *Position
	Expression Node
}

func (m *Minus) Complete(p *Position, x Node) {
	m.position = p
	m.Expression = x
}

func (m *Minus) Position() *Position { return m.position }

func (m *Minus) dump(w io.Writer, n int) {
	indent(w, n)
	fmt.Fprintf(w, "%T:%v:\n", m, m.position)
	m.Expression.dump(w, n+1)
}

type Not struct {
	position   *Position
	Expression Node
}

func (n *Not) Complete(p *Position, x Node) {
	n.position = p
	n.Expression = x
}

func (n *Not) Position() *Position { return n.position }

func (n *Not) dump(w io.Writer, lv int) {
	indent(w, lv)
	fmt.Fprintf(w, "%T:%v:\n", n, n.position)
	n.Expression.dump(w, lv+1)
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
