package renderer

import (
	"testing"

	"github.com/CephandriusMaxtori/Nativ/internal/platform"
)

func TestRegisterAndList(t *testing.T) {
	orig := backends
	backends = map[string]func() Renderer{}
	t.Cleanup(func() { backends = orig })

	list := ListBackends()
	if len(list) != 0 {
		t.Fatalf("expected empty, got %v", list)
	}

	RegisterBackend("test", func() Renderer { return nil })
	list = ListBackends()
	if len(list) != 1 || list[0] != "test" {
		t.Fatalf("expected [test], got %v", list)
	}
}

func TestNewBackend(t *testing.T) {
	orig := backends
	backends = map[string]func() Renderer{}
	t.Cleanup(func() { backends = orig })

	RegisterBackend("mock", func() Renderer {
		return &mockRenderer{}
	})

	r, err := NewBackend("mock")
	if err != nil {
		t.Fatalf("NewBackend failed: %v", err)
	}
	if r == nil {
		t.Fatal("expected non-nil renderer")
	}
}

func TestNewBackendUnknown(t *testing.T) {
	orig := backends
	backends = map[string]func() Renderer{}
	t.Cleanup(func() { backends = orig })

	_, err := NewBackend("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown backend")
	}
}

type mockRenderer struct{}

func (m *mockRenderer) Init(_ platform.Platform) error { return nil }
func (m *mockRenderer) Draw(_ []DrawCommand) error     { return nil }
func (m *mockRenderer) Shutdown() error                { return nil }

var _ Renderer = (*mockRenderer)(nil)
