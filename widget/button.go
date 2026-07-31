package widget

import (
	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

type Button struct {
	BaseWidget
	Text         string
	Font         *renderer.Font
	OnClick      func()
	hovered      bool
	pressed      bool
}

func NewButton(text string) *Button {
	return &Button{
		Text: text,
	}
}

func (b *Button) Draw(cmds *[]renderer.DrawCommand) {
	x, y, w, h := b.Bounds()
	if w <= 0 || h <= 0 {
		return
	}

	var r, g, bl, a float32
	if b.pressed {
		r, g, bl, a = 0.2, 0.5, 0.3, 1.0
	} else if b.hovered {
		r, g, bl, a = 0.3, 0.7, 0.4, 1.0
	} else {
		r, g, bl, a = 0.2, 0.6, 0.3, 1.0
	}

	*cmds = append(*cmds, renderer.DrawCommand{
		X: x, Y: y, Width: w, Height: h,
		R: r, G: g, B: bl, A: a,
	})

	// Top border
	*cmds = append(*cmds, renderer.DrawCommand{
		X: x, Y: y, Width: w, Height: 2,
		R: 0.4, G: 0.8, B: 0.5, A: 1.0,
	})

	// Button label text
	if b.Text != "" {
		textH := float32(13)
		if b.Font != nil {
			textH = float32(b.Font.Size())
		}
		*cmds = append(*cmds, renderer.DrawCommand{
			X: x + 6, Y: y + (h-textH)/2, Width: w - 12, Height: textH,
			R: 1, G: 1, B: 1, A: 1.0,
			Text: b.Text,
			Font: b.Font,
		})
	}
}

func (b *Button) HandleEvent(evt event.Event) bool {
	x, y, w, h := b.Bounds()
	mx := float32(evt.X)
	my := float32(evt.Y)

	inside := mx >= x && mx <= x+w && my >= y && my <= y+h

	switch evt.Type {
	case event.MouseMove:
		b.hovered = inside
		return inside

	case event.MouseDown:
		if inside {
			b.pressed = true
			return true
		}

	case event.MouseUp:
		if inside && b.pressed {
			if b.OnClick != nil {
				b.OnClick()
			}
			b.pressed = false
			return true
		}
		b.pressed = false
	}

	return inside
}
