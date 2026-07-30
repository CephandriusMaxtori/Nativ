package widget

import (
	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

type Widget interface {
	Bounds() (x, y, w, h float32)
	SetBounds(x, y, w, h float32)
	Draw(cmds *[]renderer.DrawCommand)
	HandleEvent(evt event.Event) bool
	Parent() Widget
	SetParent(p Widget)
}

type BaseWidget struct {
	x, y, w, h float32
	parent     Widget
}

func (b *BaseWidget) Bounds() (float32, float32, float32, float32) {
	return b.x, b.y, b.w, b.h
}

func (b *BaseWidget) SetBounds(x, y, w, h float32) {
	b.x, b.y, b.w, b.h = x, y, w, h
}

func (b *BaseWidget) Parent() Widget       { return b.parent }
func (b *BaseWidget) SetParent(p Widget)   { b.parent = p }

func (b *BaseWidget) Draw(cmds *[]renderer.DrawCommand)          {}
func (b *BaseWidget) HandleEvent(evt event.Event) bool           { return false }
