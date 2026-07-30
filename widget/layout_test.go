package widget

import (
	"testing"

	"nativ/internal/event"
	"nativ/internal/renderer"
)

func TestNewVBoxLayout(t *testing.T) {
	l := NewVBoxLayout()
	if l.dir != Vertical {
		t.Fatal("expected vertical direction")
	}
	if len(l.children) != 0 {
		t.Fatal("expected no children")
	}
}

func TestNewHBoxLayout(t *testing.T) {
	l := NewHBoxLayout()
	if l.dir != Horizontal {
		t.Fatal("expected horizontal direction")
	}
}

func TestAddChildren(t *testing.T) {
	l := NewVBoxLayout()
	b1 := NewButton("a")
	b2 := NewButton("b")

	l.Add(b1)
	if len(l.Children()) != 1 {
		t.Fatal("expected 1 child after Add")
	}
	if b1.Parent() != l {
		t.Fatal("child should have layout as parent")
	}

	l.Add(b2)
	if len(l.Children()) != 2 {
		t.Fatal("expected 2 children")
	}
}

func TestRemoveChildren(t *testing.T) {
	l := NewVBoxLayout()
	b1 := NewButton("a")
	b2 := NewButton("b")
	l.Add(b1)
	l.Add(b2)

	l.Remove(b1)
	if len(l.Children()) != 1 || l.Children()[0] != b2 {
		t.Fatal("expected only b2 remaining after remove")
	}

	l.Remove(b2)
	if len(l.Children()) != 0 {
		t.Fatal("expected no children after removing all")
	}
}

func TestRemoveNonExistentChild(t *testing.T) {
	l := NewVBoxLayout()
	l.Add(NewButton("a"))
	l.Remove(NewButton("b"))
	if len(l.Children()) != 1 {
		t.Fatal("still expected 1 child")
	}
}

func TestVerticalLayoutPositions(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)

	b1 := NewButton("top")
	b2 := NewButton("middle")
	b3 := NewButton("bottom")
	l.Add(b1)
	l.Add(b2)
	l.Add(b3)

	_, y1, _, h1 := b1.Bounds()
	_, y2, _, h2 := b2.Bounds()
	_, y3, _, h3 := b3.Bounds()

	if y1 != 10 {
		t.Fatalf("expected first child y=10, got %v", y1)
	}
	if y2 != y1+h1+8 {
		t.Fatalf("expected second child y=%v, got %v", y1+h1+8, y2)
	}
	if y3 != y2+h2+8 {
		t.Fatalf("expected third child y=%v, got %v", y2+h2+8, y3)
	}
	if h1 != h2 || h2 != h3 {
		t.Fatal("expected equal height for all children in VBox")
	}
}

func TestHorizontalLayoutPositions(t *testing.T) {
	l := NewHBoxLayout()
	l.SetBounds(0, 0, 800, 600)

	b1 := NewButton("left")
	b2 := NewButton("right")
	l.Add(b1)
	l.Add(b2)

	x1, _, _, _ := b1.Bounds()
	x2, _, w2, _ := b2.Bounds()

	if x1 != 10 {
		t.Fatalf("expected first child x=10, got %v", x1)
	}
	if x2 != x1+w2+8 {
		t.Fatalf("expected second child x=%v, got %v", x1+w2+8, x2)
	}
}

func TestLayoutDraw(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)
	l.Add(NewButton("a"))
	l.Add(NewButton("b"))

	cmds := make([]renderer.DrawCommand, 0, 10)
	l.Draw(&cmds)

	if len(cmds) != 4 {
		t.Fatalf("expected 4 commands (2 buttons x 2), got %d", len(cmds))
	}
}

func TestLayoutEventDispatch(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)
	b := NewButton("click")
	l.Add(b)

	downEvt := event.Event{Type: event.MouseDown, X: 50, Y: 15}
	if !l.HandleEvent(downEvt) {
		t.Fatal("layout should have dispatched to the button")
	}
	if !b.pressed {
		t.Fatal("button within layout should have been pressed")
	}
}

func TestLayoutEventDispatchReverseOrder(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)

	top := NewButton("top")
	bottom := NewButton("bottom")
	l.Add(top)
	l.Add(bottom)

	// Click in top button's region (y=50 is within top's bounds at ~y=10, h=286)
	evt := event.Event{Type: event.MouseDown, X: 50, Y: 50}
	l.HandleEvent(evt)

	if bottom.pressed {
		t.Fatal("bottom button should NOT be pressed (click was in top region)")
	}
	if !top.pressed {
		t.Fatal("top button SHOULD be pressed (click was in top region)")
	}
}

func TestLayoutSetBoundsTriggersRelayout(t *testing.T) {
	l := NewVBoxLayout()
	b := NewButton("test")
	l.Add(b)
	l.SetBounds(0, 0, 400, 300)

	_, y, _, h := b.Bounds()
	if y != 10 {
		t.Fatalf("expected y=10 after SetBounds, got %v", y)
	}
	if h <= 0 {
		t.Fatalf("expected positive height, got %v", h)
	}

	prevY := y
	l.SetBounds(0, 100, 400, 300)
	_, y, _, _ = b.Bounds()
	if y != prevY+100 {
		t.Fatalf("expected y=%v after moving layout down by 100, got %v", prevY+100, y)
	}
}

func TestLayoutEmptyDraw(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)
	cmds := make([]renderer.DrawCommand, 0, 10)
	l.Draw(&cmds)
	if len(cmds) != 0 {
		t.Fatal("expected no draw commands for empty layout")
	}
}

func TestLayoutNoChildren(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)
	l.SetBounds(0, 0, 0, 0)
	l.SetBounds(0, 0, 800, 600)
}

func TestLayoutAddAfterSetBounds(t *testing.T) {
	l := NewVBoxLayout()
	l.SetBounds(0, 0, 800, 600)
	b := NewButton("new")
	l.Add(b)

	_, y, _, h := b.Bounds()
	if y != 10 {
		t.Fatalf("expected y=10 for first child added after SetBounds, got %v", y)
	}
	if h <= 0 {
		t.Fatalf("expected positive height, got %v", h)
	}
}
