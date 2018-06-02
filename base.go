package mis

import (
	"bytes"
	tbox "github.com/nsf/termbox-go"
)

type Elem interface {
	Father() Elem
	SetFather(Elem)
	Son() []Elem
	Attech(Elem)
	Event(tbox.Event)
	H() EventHandler
	SetH(EventHandler)
	Pl() Place
	Drow(w, h, x, y int, bg tbox.Attribute) []Cell
}

type Place struct{ X, Y, W, H int }

func (e Place) In(x, y int) bool {
	return x >= e.X && y >= e.Y && x < e.X+e.W && y < e.Y+e.H
}

type Cell struct {
	tbox.Cell
	X, Y int
}

type BaseElemNode struct {
	S  []Elem
	f  Elem
	he EventHandler
	P  Place
}

func NewBE() BaseElemNode {
	return BaseElemNode{make([]Elem, 0), nil, nulHend{}, Place{0, 0, 0, 0}}
}

func (b *BaseElemNode) Son() []Elem {
	return b.S
}

func (b *BaseElemNode) Attech(e Elem) {
	b.S = append(b.S, e)
}

func (b *BaseElemNode) Father() Elem {
	return b.f
}

func (b *BaseElemNode) SetFather(e Elem) {
	b.f = e
}

func (b *BaseElemNode) SetH(h EventHandler) {
	b.he = h
}

func (b *BaseElemNode) H() EventHandler {
	return b.he
}

func (b *BaseElemNode) Pl() Place {
	return b.P
}

func (b *BaseElemNode) Event(e tbox.Event) {
	b.H().hend(e)
	//if bo {
	for _, el := range b.S {
		el.Event(e)
	}
	//}
}

type Color struct {
	BaseElemNode
	C  tbox.Attribute
	el Elem
}

func (c *Color) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	c.P = Place{x, y, w, h}
	cl := c.el.Drow(w, h, x, y, c.C)
	/*for i := 0; i < w; i++ {
		cl = append(cl, Cell{tbox.Cell{rune(219), c.color, c.color}, x + i, y})
		cl = append(cl, Cell{tbox.Cell{rune(219), c.color, c.color}, x + i, y + h - 1})
	}
	for i := 0; i < h; i++ {
		cl = append(cl, Cell{tbox.Cell{rune(219), c.color, c.color}, x, y + i})
		cl = append(cl, Cell{tbox.Cell{rune(219), c.color, c.color}, x + w - 1, y + i})
	}*/
	return cl
}

func C(el Elem, bg tbox.Attribute) *Color {
	return &Color{NewBE(), bg, el}
}

type nulElem struct {
	BaseElemNode
}

func (nulElem) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	cl := make([]Cell, 0)
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			cl = append(cl, Cell{tbox.Cell{rune(219), bg, bg}, i, j})
		}
	}
	return cl
}

func Fill() Elem {
	return &nulElem{NewBE()}
}

type nulElem2 struct {
	BaseElemNode
}

func (nulElem2) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	cl := make([]Cell, 0)
	return cl
}

func Nul() Elem {
	return &nulElem2{NewBE()}
}

type cC struct {
	BaseElemNode
	color tbox.Attribute
	el    Elem
}

type Frame struct {
	BaseElemNode
	P1 Place
	el Elem
}

func (f *Frame) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	f.P = Place{x, y, w, h}
	return f.el.Drow(min(w-f.P1.X, f.P1.W), min(h-f.P1.Y, f.P1.H), x+f.P1.X, y+f.P1.Y, bg)
}

func min(x, y int) int {
	if y < 0 {
		return x
	}
	if x < y {
		return x
	}
	return y
}

func F(p Place, el Elem) *Frame {
	b := NewBE()
	b.S = append(b.S, el)
	return &Frame{b, p, el}
}

type PFrame struct {
	BaseElemNode
	P1     Place
	px, py float64
	x, y   int
	el     Elem
}

func (f *PFrame) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	f.P = Place{x, y, w, h}
	f.P1.X = int(float64(w)*f.px) + x + f.x
	f.P1.Y = int(float64(h)*f.py) + y + f.y
	return f.el.Drow(min(w-f.P1.X, f.P1.W), min(h-f.P1.Y, f.P1.H), x+f.P1.X, y+f.P1.Y, bg)
}

func P(px, py float64, x, y int, w, h int, el Elem) *PFrame {
	b := NewBE()
	b.S = append(b.S, el)
	return &PFrame{b, Place{W: w, H: h}, px, py, x, y, el}
}

type d struct {
	BaseElemNode
	Drowf func(w, h, x, y int, bg tbox.Attribute) []Cell
}

func (dr *d) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	dr.P = Place{x, y, w, h}
	return dr.Drowf(w, h, x, y, bg)
}

func D(Drow func(w, h, x, y int, bg tbox.Attribute) []Cell) Elem {
	return &d{NewBE(), Drow}
}

type mList struct {
	BaseElemNode
}

func (m *mList) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	m.P = Place{x, y, w, h}
	cl := make([]Cell, 0)
	for _, el := range m.S {
		cl = append(cl, el.Drow(w, h, x, y, bg)...)
	}
	return cl
}

func MList(l ...Elem) Elem {
	be := NewBE()
	be.S = l
	return &mList{be}
}

type AnimEl struct {
	BaseElemNode
	e []Elem
	I int
}

func (a *AnimEl) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	a.P = Place{x, y, w, h}
	a.S = append([]Elem{}, a.e[a.I])
	return a.e[a.I].Drow(w, h, x, y, bg)
}

func Anim(e ...Elem) *AnimEl {
	return &AnimEl{NewBE(), e, 0}
}

type Text struct {
	BaseElemNode
	Text string
}

func NewText(txt string) *Text {
	return &Text{NewBE(), txt}
}

func (t *Text) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	t.P = Place{x, y, w, h}
	cl := make([]Cell, 0)
	i, r := 0, ' '
	for i, r = range t.Text {
		if i > w {
			break
		}
		cl = append(cl, Cell{tbox.Cell{r, 0, bg}, x + i, y})
	}
	i++
	for ; i < w; i++ {
		cl = append(cl, Cell{tbox.Cell{' ', bg, bg}, x + i, y})
	}
	return cl
}

type Paragraph struct {
	BaseElemNode
	Text *bytes.Buffer
}

func NewParagraph(txt string) *Paragraph {
	return &Paragraph{NewBE(), bytes.NewBufferString(txt)}
}

func (t *Paragraph) Drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	t.P = Place{x, y, w, h}
	cl := make([]Cell, 0)
	l, j := 0, 0
	for _, r := range t.Text.String() {
		if j > w {
			j = 0
			l++
		}
		if r == '\n' {
			j = 0
			l++
			continue
		}
		if l > h {
			break
		}
		cl = append(cl, Cell{tbox.Cell{r, 0, bg}, x + j, y + l})
		j++
	}
	return cl
}
