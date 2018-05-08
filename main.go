package main

import (
	"fmt"
	//"math"
	//"sync"
	"time"

	tbox "github.com/nsf/termbox-go"
)

func main() {
	/*i := 0
	s := Sig(50)
	m := MovArc(3, 3, 3, 30)
	for j := 0; j < 50; j++ {
		fmt.Println(s(j))
	}
	fmt.Println("-----------")
	for j := 0; j < 50; j++ {
		fmt.Println(m(s(j)))
	}
	tbox.Init()
	defer tbox.Close()
	w, h := 0, 0
	mainElem := MList()
	log := NewText("some log will appear")
	mainElem.attech(log)
	var a *AnimEl
	var f *Frame
	var once sync.Once
	f = F(
		Place{50, 4, 50, 30},
		NewCC(MList(
			Fill(),
			F(
				Place{12, 1, 7, 1},
				NewText("exemple"),
			),
			F(
				Place{3, 3, 50 - 6, -1},
				NewListOne("1", "some options"),
			),
			F(
				Place{12, 11, 6, 3},
				NewButton("send",
					func() {
						go once.Do(
							func() {
								i := 50
								for f.p1.y < h {
									time.Sleep(time.Second / 10000 * time.Duration(i))
									f.p1.y += 1
									i += 2
								}
								a.I = 1
							},
						)
					},
				),
			),
		), 235),
	)
	a = Anim(
		f,
		Nul(),
	)
	mainElem.attech(a)
	mainElem.attech(F(
		Place{5, 40, -1, -1},
		NewSwitch(),
	))
	f := F(
		Place{3, 30, 5, 3},
		NewCC(Fill(), 100),
	)
	mainElem.attech(f)
	fun := Apply(DeAcl(17), MovArc(3, 30, 40, 15), &f.p1.x, &f.p1.y)
	*/
	// --start of rend part--

}

func Rend(fps int, mainElem Elem, ctx context.Context) {
	tbox.SetInputMode( /*termbox.InputEsc |*/ tbox.InputMouse)
	tbox.SetOutputMode(tbox.Output256)
	tike = time.NewTicker(time.Second / fps)
	for {
		select {
		case <-func() chan struct{} {
			ch := make(chan struct{})
			go func() {
				evn := tbox.PollEvent()
				mainElem.event(evn)
				ch <- struct{}{}
			}()
			//log.Text = fmt.Sprint("event:", evn)
			return ch
		}():
		case <-tick.C:
		case <-ctx.Done():
			return ctx.Err()
		}
		w, h := tbox.Size()
		bg := tbox.Attribute(240)
		tbox.Clear(0, bg)
		if call && i > 40 {
			log.Text = fmt.Sprint("frame:", i)
			call = fun()
		}
		for _, c := range mainElem.drow(w, h, 0, 0, bg) {
			tbox.SetCell(c.x, c.y, c.Ch, c.Fg, c.Bg)
		}
		tbox.Flush()
	}
}

type vlist struct {
	baseElemNode
	size int
}

func (v *vlist) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	v.p = Place{x, y, w, h}
	cl := make([]Cell, 0)
	for i, el := range v.s {
		if y+i*v.size > h {
			break
		}
		cl = append(cl, el.drow(w, v.size, x, y+i*v.size, bg)...)
	}
	return cl
}

func NewVlist(size int) Elem {
	return &vlist{newBE(), size}
}

func (c *cC) drow(w, h, x, y int, bg tbox.Attribute) []Cell {
	c.p = Place{x, y, w, h}
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

func NewCC(el Elem, bg tbox.Attribute) Elem {
	b := newBE()
	b.s = append(b.s, el)
	return &cC{b, bg, el}
}

func NewButton(txt string, f func()) Elem {
	b := MList(
		Fill(),
		F(
			Place{1, 1, len(txt), 1},
			NewText(txt),
		),
	)
	a := Anim(
		NewCC(b, 242),
		NewCC(b, 238),
	)

	g := func(x, y int, k tbox.Key, drag bool) {
		a.I = 1
		if k == tbox.MouseLeft && !drag {
			f()
		}
		go func() {
			time.Sleep(time.Second / 10)
			a.I = 0
		}()
	}
	return F(Place{x: 0, y: 0, w: len(txt) + 2, h: 3}, OnMouse(a, g))
}

func NewList(a ...string) Elem {
	l := NewVlist(4)
	for _, txt := range a {
		an := Anim(
			C(NewText("\u25cb"), 237),
			C(NewText("\u25c9"), 237 /*3*/),
		)
		an.I = 0
		t := NewText(txt)
		l.attech(
			F(
				Place{0, 0, -1, 3},
				OnMouse(
					C(
						MList(
							Fill(),
							F(
								Place{4, 1, -1, 1},
								t,
							),
							F(
								Place{1, 1, 2, 1},
								an,
							),
						),
						234,
					),
					func(x, y int, k tbox.Key, drag bool) {
						if k == tbox.MouseLeft && !drag {
							an.I = 1 - an.I
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
	Elem
}

func NewSwitch() switchEl {
	sEl := switchEl{}
	f := F(
		Place{0, 0 /*will chagne*/, 3, 3},
		D(func(w, h, x, y int, bg tbox.Attribute) []Cell {
			cl := make([]Cell, 0)
			cl = append(cl, Cell{tbox.Cell{' ', 100, 100}, x + 1, y + 1})

			cl = append(cl, Cell{tbox.Cell{'\u2581', 100, bg}, x + 1, y})
			cl = append(cl, Cell{tbox.Cell{'\u2587', bg, 100}, x + 1, y + 2})
			//cl = append(cl, Cell{tbox.Cell{'\u258a', bg, 100}, x, y + 1})
			//cl = append(cl, Cell{tbox.Cell{'\u258d', 100, bg}, x + 2, y + 1})
			return cl
		}),
	)
	c := C(Fill(), 189)
	s := MList(
		F(
			Place{1, 1, 3, 1},
			c,
		),
		f,
	)
	sMouse := F(
		Place{0, 0, 5, 3},
		OnMouse(
			s,
			func(x, y int, k tbox.Key, drag bool) {
				go func() {
					sEl.B = !sEl.B
					if k == tbox.MouseLeft && !drag {
						f.p1.x = 2 - f.p1.x
					}
				}()
			},
		),
	)
	sEl.Elem = sMouse
	sEl.f = 0.0
	sEl.B = false
	return sEl
}

type FoucosGroup struct {
	f func()
	i int
}

func (f *FoucosGroup) F(i int, out func()) {
	f.i = i
	f.f()
	f.f = out
	return
}

func (f *FoucosGroup) Get(i int) bool {
	return f.i == i
}

func NewListOne(a ...string) Elem {
	l := NewVlist(4)
	f := new(FoucosGroup)
	f.f = func() {}
	for j, txt := range a {
		j1 := j
		an := Anim(
			C(NewText("\u25cb"), 237),
			C(NewText("\u25c9"), 237 /*3*/),
		)
		an.I = 0
		t := NewText(txt)
		l.attech(
			F(
				Place{0, 0, -1, 3},
				OnMouse(
					C(
						MList(
							Fill(),
							F(
								Place{4, 1, -1, 1},
								t,
							),
							F(
								Place{1, 1, 2, 1},
								an,
							),
						),
						234,
					),
					func(x, y int, k tbox.Key, drag bool) {
						if k == tbox.MouseLeft && !drag {
							if !f.Get(j1) {
								f.F(j1, func() {
									an.I = 0
								},
								)
							}
							an.I = 1
						}
					},
				),
			),
		)
	}
	return l
}
