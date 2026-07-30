package widget

import (
	"testing"

	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

func TestBaseWidgetBounds(t *testing.T) {
	w := &BaseWidget{}
	x, y, bw, bh := w.Bounds()
	if x != 0 || y != 0 || bw != 0 || bh != 0 {
		t.Fatalf("expected zero bounds, got %v %v %v %v", x, y, bw, bh)
	}

	w.SetBounds(10, 20, 100, 200)
	x, y, bw, bh = w.Bounds()
	if x != 10 || y != 20 || bw != 100 || bh != 200 {
		t.Fatalf("expected (10,20,100,200), got (%v,%v,%v,%v)", x, y, bw, bh)
	}
}

func TestBaseWidgetParent(t *testing.T) {
	child := &BaseWidget{}
	parent := &BaseWidget{}

	if child.Parent() != nil {
		t.Fatal("expected nil parent")
	}

	child.SetParent(parent)
	if child.Parent() != parent {
		t.Fatal("expected parent to be set")
	}
}

func TestBaseWidgetDrawNoop(t *testing.T) {
	w := &BaseWidget{}
	cmds := make([]renderer.DrawCommand, 0, 10)
	w.Draw(&cmds)
	if len(cmds) != 0 {
		t.Fatal("BaseWidget.Draw should not produce commands")
	}
}

func TestBaseWidgetHandleEvent(t *testing.T) {
	w := &BaseWidget{}
	evt := event.Event{Type: event.MouseDown, X: 0, Y: 0}
	if w.HandleEvent(evt) {
		t.Fatal("BaseWidget.HandleEvent should return false")
	}
}
