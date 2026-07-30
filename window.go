package nativ

import (
	"fmt"

	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/platform"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
	"github.com/CephandriusMaxtori/Nativ/widget"
)

type Window struct {
	platform platform.Platform
	renderer renderer.Renderer
	root     widget.Widget
	title    string
	width    int
	height   int
}

func NewWindow(title string, width, height int, rendererName string) (*Window, error) {
	p, err := platform.New()
	if err != nil {
		return nil, fmt.Errorf("nativ: %w", err)
	}

	_, _, err = p.CreateWindow(title, width, height)
	if err != nil {
		return nil, fmt.Errorf("nativ: %w", err)
	}

	r, err := renderer.NewBackend(rendererName)
	if err != nil {
		return nil, fmt.Errorf("nativ: %w", err)
	}

	if err := r.Init(p); err != nil {
		return nil, fmt.Errorf("nativ: %w", err)
	}

	return &Window{
		platform: p,
		renderer: r,
		title:    title,
		width:    width,
		height:   height,
	}, nil
}

func (w *Window) SetRoot(root widget.Widget) {
	if root == nil {
		return
	}
	w.root = root
	w.root.SetBounds(0, 0, float32(w.width), float32(w.height))
}

func (w *Window) Root() widget.Widget {
	return w.root
}

func (w *Window) Show() {
	// Window is shown during Run()
}

func (w *Window) Run() error {
	if w.root != nil {
		w.root.SetBounds(0, 0, float32(w.width), float32(w.height))
	}

	return w.platform.Pump(
		func(evt event.Event) {
			w.handleEvent(evt)
		},
		func() {
			w.renderFrame()
		},
	)
}

func (w *Window) Close() {
	w.platform.PostQuit()
}

func (w *Window) handleEvent(evt event.Event) {
	switch evt.Type {
	case event.Close:
		w.Close()
		return

	case event.Resize:
		w.width = evt.Width
		w.height = evt.Height
		if w.root != nil {
			w.root.SetBounds(0, 0, float32(w.width), float32(w.height))
		}
	}

	if w.root != nil {
		w.root.HandleEvent(evt)
	}
}

func (w *Window) renderFrame() {
	cmds := make([]renderer.DrawCommand, 0, 1024)
	if w.root != nil {
		w.root.Draw(&cmds)
	}
	w.renderer.Draw(cmds)
}
