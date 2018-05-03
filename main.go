package main

import (
	//"fmt"
	"time"

	tbox "github.com/nsf/termbox-go"
)

func main() {
	i := 0
	tbox.Init()
	defer tbox.Close()
	mainElem := NewMList()
	//log := NewText("some log will appear")
	//mainElem.attech(log)
	var a *anim
	a = NewAnim(
		F(
			elSize{50, 4, 50, 30},
			NewCC(NewMList(
				NewNul(),
				F(
					elSize{12, 1, 7, 1},
					NewText("exemple"),
				),
				F(
					elSize{3, 3, 50 - 6, -1},
					NewList("1", "some options"),
				),
				F(
					elSize{12, 11, 6, 3},
					NewButton("send", func() {
						go func() {
							time.Sleep(time.Second / 7)
							a.i = 1
						}()
					}),
				),
			), 235),
		),
		NewNul2(),
	).(*anim)
	mainElem.attech(a)
	mainElem.attech(F(
		elSize{5, 40, -1, -1},
		NewSwitch(),
	))

	// --start of rend part--

	tbox.SetInputMode( /*termbox.InputEsc |*/ tbox.InputMouse)
	tbox.SetOutputMode(tbox.Output256)
	ev := make(chan struct{}, 0)
	//catch events
	go func() {
		for {
			evn := tbox.PollEvent()
			ev <- struct{}{}
			//log.Text = fmt.Sprint("event:", evn)
			mainElem.event(evn)
		}
	}()
	//the rend funcion
	for i = 0; i < 50; i++ {
		ch := time.After(time.Second / 10)
		w, h := tbox.Size()
		//log.Text = fmt.Sprint("size:", w, h)
		bg := tbox.Attribute(240)
		tbox.Clear(0, bg)
		for _, c := range mainElem.drow(w, h, 0, 0, bg) {
			tbox.SetCell(c.x, c.y, c.Ch, c.Fg, c.Bg)
		}
		tbox.Flush()
		select {
		case <-ev:
		case <-ch:
		}
	}
}

type elem interface {
	father() elem
	setFather(elem)
	son() []elem
	attech(elem)
	event(tbox.Event)
	h() EventHandler
	setH(EventHandler)
	pl() elSize
	drow(w, h, x, y int, bg tbox.Attribute) []Cell
}

type elSize struct{ x, y, w, h int }

func (e elSize) in(x, y int) bool {
	return x >= e.x && y >= e.y && x < e.x+e.w && y < e.y+e.h
}

type Cell struct {
	tbox.Cell
	x, y int
}

type baseElemNode struct {
	s  []elem
	f  elem
	d  elem
	u  elem
	he EventHandler
	p  elSize
}

//return true if need to pass on
type EventHandler interface {
	hend(tbox.Event) bool
}

type nulHend struct{}

func (nulHend) hend(tbox.Event) bool {
	return true
}

type splitHend struct {
	mouse func(x, y int, k tbox.Key)
	key   func(k tbox.Key, Mod tbox.Modifier, ch rune)
}

func (s splitHend) hend(e tbox.Event) bool {
	if e.Type == tbox.EventMouse {
		s.mouse(e.MouseX, e.MouseY, e.Key)
	} else if e.Type == tbox.EventKey {
		s.key(e.Key, e.Mod, e.Ch)
	}
	return true
}

func OnMouse(b elem, h func(x, y int, k tbox.Key)) elem {
	b.setH(splitHend{
		func(x, y int, k tbox.Key) {
			if b.pl().in(x, y) {
				h(x, y, k)
			}
		},
		func(k tbox.Key, Mod tbox.Modifier, ch rune) {

		},
	})
	return b
}

func newBE() baseElemNode {
	return baseElemNode{make([]elem, 0), nil, nil, nil, nulHend{}, elSize{0, 0, 0, 0}}
}

func (b *baseElemNode) son() []elem {
	return b.s
}

func (b *baseElemNode) attech(e elem) {
	b.s = append(b.s, e)
}

func (b *baseElemNode) father() elem {
	return b.f
}

func (b *baseElemNode) setFather(e elem) {
	b.f = e
}

func (b *baseElemNode) setH(h EventHandler) {
	b.he = h
}

func (b *baseElemNode) h() EventHandler {
	return b.he
}

func (b *baseElemNode) pl() elSize {
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

type vlist struct {
	baseElemNode
	size int
}

func (v *vlist) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	v.p = elSize{x, y, w, h}
	cl := make([]Cell, 0)
	for i, el := range v.s {
		if y+i*v.size > h {
			break
		}
		cl = append(cl, el.drow(w, v.size, x, y+i*v.size, bg)...)
	}
	return cl
}

func NewVlist(size int) elem {
	return &vlist{newBE(), size}
}

type Text struct {
	baseElemNode
	Text string
}

func NewText(txt string) *Text {
	return &Text{newBE(), txt}
}

func (t *Text) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	t.p = elSize{x, y, w, h}
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

type colorFrame struct {
	baseElemNode
	color tbox.Attribute
	el    elem
}

func (c *colorFrame) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	c.p = elSize{x, y, w, h}
	cl := c.el.drow(w, h, x, y, c.color)
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

func NewColor(el elem, bg tbox.Attribute) elem {
	return &colorFrame{newBE(), bg, el}
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

func NewNul() elem {
	return &nulElem{newBE()}
}

type nulElem2 struct {
	baseElemNode
}

func (nulElem2) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	cl := make([]Cell, 0)
	return cl
}

func NewNul2() elem {
	return &nulElem2{newBE()}
}

type cC struct {
	baseElemNode
	color tbox.Attribute
	el    elem
}

func (c *cC) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	c.p = elSize{x, y, w, h}
	if !(w > 1 && h > 1) {
		return nil
	}
	cl := c.el.drow(w, h, x, y, c.color)
	for i := 0; i < w-1; i++ {
		cl = append(cl, Cell{tbox.Cell{rune(9608), c.color, c.color}, x + i, y})
		cl = append(cl, Cell{tbox.Cell{rune(9608), c.color, c.color}, x + i, y + h - 1})
	}
	for i := 0; i < h-1; i++ {
		cl = append(cl, Cell{tbox.Cell{rune(9608), c.color, c.color}, x, y + i})
		cl = append(cl, Cell{tbox.Cell{rune(9608), c.color, c.color}, x + w - 1, y + i})
	}
	cl = append(cl,
		Cell{tbox.Cell{rune(9600), bg, c.color}, x, y},
		Cell{tbox.Cell{rune(9600), bg, c.color}, x + w - 1, y},
		Cell{tbox.Cell{rune(9604), bg, c.color}, x, y + h - 1},
		Cell{tbox.Cell{rune(9604), bg, c.color}, x + w - 1, y + h - 1},
	)
	return cl
}

func NewCC(el elem, bg tbox.Attribute) elem {
	b := newBE()
	b.s = append(b.s, el)
	return &cC{b, bg, el}
}

type frame struct {
	baseElemNode
	p1 elSize
	el elem
}

func (f *frame) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	f.p = elSize{x, y, w, h}
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

func F(p elSize, el elem) elem {
	b := newBE()
	b.s = append(b.s, el)
	return &frame{b, p, el}
}

type d struct {
	baseElemNode
	drowf func(w, h, x, y int, bg tbox.Attribute) []Cell
}

func (dr *d) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	dr.p = elSize{x, y, w, h}
	return dr.drowf(w, h, x, y, bg)
}

func D(drow func(w, h, x, y int, bg tbox.Attribute) []Cell) elem {
	return &d{newBE(), drow}
}

type mList struct {
	baseElemNode
}

func (m *mList) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	m.p = elSize{x, y, w, h}
	cl := make([]Cell, 0)
	for _, el := range m.s {
		cl = append(cl, el.drow(w, h, x, y, bg)...)
	}
	return cl
}

func NewMList(l ...elem) elem {
	be := newBE()
	be.s = l
	return &mList{be}
}

type anim struct {
	baseElemNode
	e []elem
	i int
}

func (a *anim) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	a.p = elSize{x, y, w, h}
	a.s = append([]elem{}, a.e[a.i])
	return a.e[a.i].drow(w, h, x, y, bg)
}

func NewAnim(e ...elem) elem {
	return &anim{newBE(), e, 0}
}

func NewButton(txt string, f func()) elem {
	b := NewMList(
		NewNul(),
		F(
			elSize{1, 1, len(txt), 1},
			NewText(txt),
		),
	)
	a := NewAnim(
		NewCC(b, 242),
		NewCC(b, 238),
	).(*anim)

	g := func(x, y int, k tbox.Key) {
		a.i = 1
		f()
		go func() {
			time.Sleep(time.Second / 10)
			a.i = 0
		}()
	}
	return F(elSize{x: 0, y: 0, w: len(txt) + 2, h: 3}, OnMouse(a, g))
}

//NewList create new list
//List := def
//VList{
//	Mlist{
//		F{
//			OnMouse{
//				anim{
//					Color(text(" ")),
//					Color(text("V"))
//				}
//			}
//		},
//		F{
//			Text()
//		}
//	}
//}
func NewList(a ...string) elem {
	l := NewVlist(4)
	for _, txt := range a {
		an := NewAnim(
			NewColor(NewText("\u25cb"), 237),
			NewColor(NewText("\u25c9"), 237 /*3*/),
		).(*anim)
		an.i = 0
		t := NewText(txt)
		l.attech(
			F(
				elSize{0, 0, -1, 3},
				OnMouse(
					NewColor(
						NewMList(
							NewNul(),
							F(
								elSize{4, 1, -1, 1},
								t,
							),
							F(
								elSize{1, 1, 2, 1},
								an,
							),
						),
						234,
					),
					func(x, y int, k tbox.Key) {
						if k == tbox.MouseLeft {
							an.i = 1 - an.i
						}
					},
				),
			),
		)
	}
	return l
}

type switchEl struct {
	B bool
	f float64
	elem
}

func NewSwitch() switchEl {
	f := F(
		elSize{0, 0 /*will chagne*/, 3, 3},
		D(func(w, h, x, y int, bg tbox.Attribute) []Cell {
			cl := make([]Cell, 0)
			cl = append(cl, Cell{tbox.Cell{' ', 100, 100}, x + 1, y + 1})

			cl = append(cl, Cell{tbox.Cell{'\u2581', 100, bg}, x + 1, y})
			cl = append(cl, Cell{tbox.Cell{'\u2587', bg, 100}, x + 1, y + 2})
			//cl = append(cl, Cell{tbox.Cell{'\u258a', bg, 100}, x, y + 1})
			//cl = append(cl, Cell{tbox.Cell{'\u258d', 100, bg}, x + 2, y + 1})
			return cl
		}),
	).(*frame)
	c := NewColor(NewNul(), 189).(*colorFrame)
	s := NewMList(
		F(
			elSize{1, 1, 3, 1},
			c,
		),
		f,
	)
	sMouse := F(
		elSize{0, 0, 5, 3},
		OnMouse(
			s,
			func(x, y int, k tbox.Key) {
				go func() {
					if k == tbox.MouseLeft {
						f.p1.x = 2 - f.p1.x
					}
				}()
			},
		),
	)
	return switchEl{false, 0.0, sMouse}
}
