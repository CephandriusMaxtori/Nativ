package widget

import (
	"nativ/internal/event"
	"nativ/internal/renderer"
)

type Button struct {
	BaseWidget
	Text         string
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
		r, g, bl, a = 0.2, 0.5, 0.3, 1.0 // dark green when pressed
	} else if b.hovered {
		r, g, bl, a = 0.3, 0.7, 0.4, 1.0 // lighter green when hovered
	} else {
		r, g, bl, a = 0.2, 0.6, 0.3, 1.0 // normal green
	}

	*cmds = append(*cmds, renderer.DrawCommand{
		X: x, Y: y, Width: w, Height: h,
		R: r, G: g, B: bl, A: a,
	})

	// Draw a thin border
	borderColor := renderer.DrawCommand{
		X: x, Y: y, Width: w, Height: 2,
		R: 0.4, G: 0.8, B: 0.5, A: 1.0,
	}
	*cmds = append(*cmds, borderColor)
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
