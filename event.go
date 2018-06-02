package mis

import (
	tbox "github.com/nsf/termbox-go"
)

const (
	EventUpdate tbox.EventType = tbox.EventNone + 1
)

//return true if need to pass on
type EventHandler interface {
	hend(tbox.Event) bool
}

type nulHend struct{}

func (nulHend) hend(tbox.Event) bool {
	return true
}

type splitHend struct {
	mouse  func(x, y int, k tbox.Key, drag bool)
	key    func(k tbox.Key, Mod tbox.Modifier, ch rune)
	update func()
}

func (s splitHend) hend(e tbox.Event) bool {
	if e.Type == EventUpdate {
		s.update()
	} else if e.Type == tbox.EventMouse {
		s.mouse(e.MouseX, e.MouseY, e.Key, e.Mod == 2)
	} else if e.Type == tbox.EventKey {
		s.key(e.Key, e.Mod, e.Ch)
	}
	return true
}

func OnMouse(b Elem, h func(x, y int, k tbox.Key, d bool)) Elem {
	if s, ok := b.H().(splitHend); ok {
		s.mouse = func(x, y int, k tbox.Key, d bool) {
			if b.Pl().In(x, y) {
				h(x, y, k, d)
			}
		}
		b.SetH(s)
	}
	b.SetH(splitHend{
		func(x, y int, k tbox.Key, d bool) {
			if b.Pl().In(x, y) {
				h(x, y, k, d)
			}
		},
		func(k tbox.Key, Mod tbox.Modifier, ch rune) {

		},
		func() {

		},
	})
	return b
}

func OnKey(b Elem, h func(k tbox.Key, Mod tbox.Modifier, ch rune)) Elem {
	if s, ok := b.H().(splitHend); ok {
		s.key = h
		b.SetH(s)
	}
	b.SetH(splitHend{
		func(x, y int, k tbox.Key, d bool) {

		},
		h,
		func() {

		},
	})
	return b
}

func OnUpdate(b Elem, h func()) Elem {
	if s, ok := b.H().(splitHend); ok {
		s.update = h
		b.SetH(s)
	}
	b.SetH(splitHend{
		func(x, y int, k tbox.Key, d bool) {

		},
		func(k tbox.Key, Mod tbox.Modifier, ch rune) {

		},

		h,
	})
	return b
}
