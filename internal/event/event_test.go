package event

import "testing"

func TestEventFields(t *testing.T) {
	e := Event{Type: MouseDown, X: 10, Y: 20}
	if e.Type != MouseDown || e.X != 10 || e.Y != 20 {
		t.Fatalf("unexpected event fields")
	}
}

func TestEventResize(t *testing.T) {
	e := Event{Type: Resize, Width: 800, Height: 600}
	if e.Width != 800 || e.Height != 600 {
		t.Fatalf("unexpected resize fields")
	}
}

func TestEventClose(t *testing.T) {
	e := Event{Type: Close}
	if e.Type != Close {
		t.Fatalf("expected Close event type")
	}
}
