package widget

import (
	"nativ/internal/event"
	"nativ/internal/renderer"
)

type Label struct {
	BaseWidget
	Text string
}

func NewLabel(text string) *Label {
	return &Label{
		Text: text,
	}
}

func (l *Label) Draw(cmds *[]renderer.DrawCommand) {
	x, y, w, h := l.Bounds()
	if w <= 0 || h <= 0 {
		return
	}

	*cmds = append(*cmds, renderer.DrawCommand{
		X: x, Y: y, Width: w, Height: h,
		R: 0.15, G: 0.15, B: 0.18, A: 1.0,
	})
}

func (l *Label) HandleEvent(evt event.Event) bool {
	x, y, w, h := l.Bounds()
	mx := float32(evt.X)
	my := float32(evt.Y)

	inside := mx >= x && mx <= x+w && my >= y && my <= y+h
	return inside
}
