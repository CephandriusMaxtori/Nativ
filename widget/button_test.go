package widget

import (
	"testing"

	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

func TestNewButton(t *testing.T) {
	b := NewButton("test")
	if b.Text != "test" {
		t.Fatalf("expected text 'test', got %q", b.Text)
	}
	if b.OnClick != nil {
		t.Fatal("expected nil OnClick")
	}
}

func TestButtonDraw(t *testing.T) {
	b := NewButton("test")
	b.SetBounds(0, 0, 100, 30)

	cmds := make([]renderer.DrawCommand, 0, 10)
	b.Draw(&cmds)

	if len(cmds) == 0 {
		t.Fatal("Button.Draw should produce at least one command")
	}

	body := cmds[0]
	if body.X != 0 || body.Y != 0 || body.Width != 100 || body.Height != 30 {
		t.Fatalf("unexpected body bounds: %v %v %v %v", body.X, body.Y, body.Width, body.Height)
	}
}

func TestButtonDrawZeroBounds(t *testing.T) {
	b := NewButton("test")
	cmds := make([]renderer.DrawCommand, 0, 10)
	b.Draw(&cmds)
	if len(cmds) != 0 {
		t.Fatal("Button.Draw should produce no commands for zero bounds")
	}
}

func TestButtonHoverEnterLeave(t *testing.T) {
	b := NewButton("test")
	b.SetBounds(0, 0, 100, 30)

	evt := event.Event{Type: event.MouseMove, X: 50, Y: 15}
	handled := b.HandleEvent(evt)
	if !handled {
		t.Fatal("expected hover inside to return true")
	}
	if !b.hovered {
		t.Fatal("expected hovered=true after move inside")
	}

	evt = event.Event{Type: event.MouseMove, X: 200, Y: 200}
	handled = b.HandleEvent(evt)
	if handled {
		t.Fatal("expected hover outside to return false")
	}
	if b.hovered {
		t.Fatal("expected hovered=false after move outside")
	}
}

func TestButtonClick(t *testing.T) {
	b := NewButton("test")
	b.SetBounds(0, 0, 100, 30)

	clicked := false
	b.OnClick = func() { clicked = true }

	downEvt := event.Event{Type: event.MouseDown, X: 50, Y: 15}
	if !b.HandleEvent(downEvt) {
		t.Fatal("expected MouseDown inside to return true")
	}
	if !b.pressed {
		t.Fatal("expected pressed=true after MouseDown inside")
	}

	upEvt := event.Event{Type: event.MouseUp, X: 50, Y: 15}
	if !b.HandleEvent(upEvt) {
		t.Fatal("expected MouseUp inside to return true")
	}
	if !clicked {
		t.Fatal("expected OnClick to have been called")
	}
	if b.pressed {
		t.Fatal("expected pressed=false after MouseUp")
	}
}

func TestButtonClickOutside(t *testing.T) {
	b := NewButton("test")
	b.SetBounds(0, 0, 100, 30)

	clicked := false
	b.OnClick = func() { clicked = true }

	b.HandleEvent(event.Event{Type: event.MouseDown, X: 50, Y: 15})
	b.HandleEvent(event.Event{Type: event.MouseUp, X: 200, Y: 200})
	if clicked {
		t.Fatal("expected OnClick NOT to have been called for release outside")
	}
}

func TestButtonClickedFlagOnMouseUpOnly(t *testing.T) {
	b := NewButton("test")
	b.SetBounds(0, 0, 100, 30)

	clicked := false
	b.OnClick = func() { clicked = true }

	b.HandleEvent(event.Event{Type: event.MouseDown, X: 50, Y: 15})
	if clicked {
		t.Fatal("OnClick should not fire on MouseDown")
	}

	b.HandleEvent(event.Event{Type: event.MouseDown, X: 50, Y: 15})
	if clicked {
		t.Fatal("OnClick should not fire on second MouseDown")
	}
}
