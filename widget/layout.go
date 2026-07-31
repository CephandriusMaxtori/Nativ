package widget

import (
	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

type Direction uint8

const (
	Vertical Direction = iota
	Horizontal
)

type BoxLayout struct {
	BaseWidget
	children []Widget
	dir      Direction
	padding  float32
	spacing  float32
}

func NewVBoxLayout() *BoxLayout {
	return &BoxLayout{
		dir:     Vertical,
		padding: 10,
		spacing: 8,
	}
}

func NewHBoxLayout() *BoxLayout {
	return &BoxLayout{
		dir:     Horizontal,
		padding: 10,
		spacing: 8,
	}
}

func (l *BoxLayout) Add(w Widget) {
	w.SetParent(l)
	l.children = append(l.children, w)
	l.relayout()
}

func (l *BoxLayout) Remove(w Widget) {
	for i, c := range l.children {
		if c == w {
			l.children = append(l.children[:i], l.children[i+1:]...)
			break
		}
	}
	l.relayout()
}

func (l *BoxLayout) Children() []Widget {
	return l.children
}

func (l *BoxLayout) relayout() {
	px, py, pw, ph := l.Bounds()
	if pw <= 0 || ph <= 0 {
		pw = 800
		ph = 600
	}
	if l.padding*2 >= ph {
		l.padding = ph / 4
	}
	if l.padding*2 >= pw {
		l.padding = pw / 4
	}
	if l.spacing >= ph {
		l.spacing = ph / 4
	}

	if l.dir == Vertical {
		totalContentH := max(ph-2*l.padding, 0)
		contentW := max(pw-2*l.padding, 0)
		childCount := len(l.children)

		if childCount == 0 {
			return
		}

		childH := (totalContentH - float32(childCount-1)*l.spacing) / float32(childCount)
		if childH < 20 {
			childH = 30
		}
		childH = min(childH, ph)

		y := py + l.padding
		for _, c := range l.children {
			c.SetBounds(px+l.padding, y, contentW, childH)
			y += childH + l.spacing
		}
	} else {
		totalContentW := max(pw-2*l.padding, 0)
		contentH := max(ph-2*l.padding, 0)
		childCount := len(l.children)

		if childCount == 0 {
			return
		}

		childW := (totalContentW - float32(childCount-1)*l.spacing) / float32(childCount)
		if childW < 40 {
			childW = 80
		}
		childW = min(childW, pw)

		x := px + l.padding
		for _, c := range l.children {
			c.SetBounds(x, py+l.padding, childW, contentH)
			x += childW + l.spacing
		}
	}
}

func (l *BoxLayout) SetBounds(x, y, w, h float32) {
	l.BaseWidget.SetBounds(x, y, w, h)
	l.relayout()
}

func (l *BoxLayout) Draw(cmds *[]renderer.DrawCommand) {
	for _, c := range l.children {
		c.Draw(cmds)
	}
}

func (l *BoxLayout) HandleEvent(evt event.Event) bool {
	// Dispatch to children in reverse order (topmost first)
	for i := len(l.children) - 1; i >= 0; i-- {
		if l.children[i].HandleEvent(evt) {
			return true
		}
	}
	return false
}
