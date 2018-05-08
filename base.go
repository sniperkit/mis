package main

import (
	tbox "github.com/nsf/termbox-go"
)

type Elem interface {
	father() Elem
	setFather(Elem)
	son() []Elem
	attech(Elem)
	event(tbox.Event)
	h() EventHandler
	setH(EventHandler)
	pl() Place
	drow(w, h, x, y int, bg tbox.Attribute) []Cell
}

type Place struct{ x, y, w, h int }

func (e Place) in(x, y int) bool {
	return x >= e.x && y >= e.y && x < e.x+e.w && y < e.y+e.h
}

type Cell struct {
	tbox.Cell
	x, y int
}

type baseElemNode struct {
	s  []Elem
	f  Elem
	he EventHandler
	p  Place
}

func newBE() baseElemNode {
	return baseElemNode{make([]Elem, 0), nil, nulHend{}, Place{0, 0, 0, 0}}
}

func (b *baseElemNode) son() []Elem {
	return b.s
}

func (b *baseElemNode) attech(e Elem) {
	b.s = append(b.s, e)
}

func (b *baseElemNode) father() Elem {
	return b.f
}

func (b *baseElemNode) setFather(e Elem) {
	b.f = e
}

func (b *baseElemNode) setH(h EventHandler) {
	b.he = h
}

func (b *baseElemNode) h() EventHandler {
	return b.he
}

func (b *baseElemNode) pl() Place {
	return b.p
}

func (b *baseElemNode) event(e tbox.Event) {
	b.h().hend(e)
	//if bo {
	for _, el := range b.s {
		el.event(e)
	}
	//}
}

type Color struct {
	baseElemNode
	C  tbox.Attribute
	el Elem
}

func (c *Color) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	c.p = Place{x, y, w, h}
	cl := c.el.drow(w, h, x, y, c.C)
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
	return &Color{newBE(), bg, el}
}

type nulElem struct {
	baseElemNode
}

func (nulElem) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	cl := make([]Cell, 0)
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			cl = append(cl, Cell{tbox.Cell{rune(219), bg, bg}, i, j})
		}
	}
	return cl
}

func Fill() Elem {
	return &nulElem{newBE()}
}

type nulElem2 struct {
	baseElemNode
}

func (nulElem2) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	cl := make([]Cell, 0)
	return cl
}

func Nul() Elem {
	return &nulElem2{newBE()}
}

type cC struct {
	baseElemNode
	color tbox.Attribute
	el    Elem
}

type Frame struct {
	baseElemNode
	p1 Place
	el Elem
}

func (f *Frame) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	f.p = Place{x, y, w, h}
	return f.el.drow(min(w-f.p1.x, f.p1.w), min(h-f.p1.y, f.p1.h), x+f.p1.x, y+f.p1.y, bg)
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
	b := newBE()
	b.s = append(b.s, el)
	return &Frame{b, p, el}
}

type d struct {
	baseElemNode
	drowf func(w, h, x, y int, bg tbox.Attribute) []Cell
}

func (dr *d) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	dr.p = Place{x, y, w, h}
	return dr.drowf(w, h, x, y, bg)
}

func D(drow func(w, h, x, y int, bg tbox.Attribute) []Cell) Elem {
	return &d{newBE(), drow}
}

type mList struct {
	baseElemNode
}

func (m *mList) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	m.p = Place{x, y, w, h}
	cl := make([]Cell, 0)
	for _, el := range m.s {
		cl = append(cl, el.drow(w, h, x, y, bg)...)
	}
	return cl
}

func MList(l ...Elem) Elem {
	be := newBE()
	be.s = l
	return &mList{be}
}

type AnimEl struct {
	baseElemNode
	e []Elem
	I int
}

func (a *AnimEl) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	a.p = Place{x, y, w, h}
	a.s = append([]Elem{}, a.e[a.I])
	return a.e[a.I].drow(w, h, x, y, bg)
}

func Anim(e ...Elem) *AnimEl {
	return &AnimEl{newBE(), e, 0}
}

type Text struct {
	baseElemNode
	Text string
}

func NewText(txt string) *Text {
	return &Text{newBE(), txt}
}

func (t *Text) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	t.p = Place{x, y, w, h}
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
