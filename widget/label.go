package widget

import (
	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

type Label struct {
	BaseWidget
	Text string
	Font *renderer.Font
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

	// Background
	*cmds = append(*cmds, renderer.DrawCommand{
		X: x, Y: y, Width: w, Height: h,
		R: 0.15, G: 0.15, B: 0.18, A: 1.0,
	})

	// Text label
	if l.Text != "" {
		*cmds = append(*cmds, renderer.DrawCommand{
			X: x + 4, Y: y + 4, Width: w - 8, Height: h - 8,
			R: 0.9, G: 0.9, B: 0.95, A: 1.0,
			Text: l.Text,
			Font: l.Font,
		})
	}
}

func (l *Label) HandleEvent(evt event.Event) bool {
	x, y, w, h := l.Bounds()
	mx := float32(evt.X)
	my := float32(evt.Y)

	inside := mx >= x && mx <= x+w && my >= y && my <= y+h
	return inside
}
