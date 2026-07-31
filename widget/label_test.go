package widget

import (
	"testing"

	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

func TestNewLabel(t *testing.T) {
	l := NewLabel("hello")
	if l.Text != "hello" {
		t.Fatalf("expected text 'hello', got %q", l.Text)
	}
}

func TestLabelDraw(t *testing.T) {
	l := NewLabel("hello")
	l.SetBounds(10, 20, 200, 30)

	cmds := make([]renderer.DrawCommand, 0, 10)
	l.Draw(&cmds)

	if len(cmds) != 2 {
		t.Fatalf("expected 2 commands (background + text), got %d", len(cmds))
	}

	c := cmds[0]
	if c.X != 10 || c.Y != 20 || c.Width != 200 || c.Height != 30 {
		t.Fatalf("unexpected bounds: %v %v %v %v", c.X, c.Y, c.Width, c.Height)
	}
	if cmds[1].Text != "hello" {
		t.Fatalf("expected text command, got %+v", cmds[1])
	}
}

func TestLabelDrawZeroBounds(t *testing.T) {
	l := NewLabel("hello")
	cmds := make([]renderer.DrawCommand, 0, 10)
	l.Draw(&cmds)
	if len(cmds) != 0 {
		t.Fatal("Label.Draw should produce no commands for zero bounds")
	}
}

func TestLabelHitTest(t *testing.T) {
	l := NewLabel("hello")
	l.SetBounds(0, 0, 100, 30)

	inside := l.HandleEvent(event.Event{Type: event.MouseMove, X: 50, Y: 15})
	if !inside {
		t.Fatal("expected hit inside")
	}

	outside := l.HandleEvent(event.Event{Type: event.MouseMove, X: 200, Y: 50})
	if outside {
		t.Fatal("expected no hit outside")
	}
}
