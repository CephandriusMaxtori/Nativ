package renderer

import (
	"errors"
	"fmt"

	"nativ/internal/platform"
)

type DrawCommand struct {
	X, Y, Width, Height float32
	R, G, B, A          float32
}

type Renderer interface {
	Init(platform.Platform) error
	Draw([]DrawCommand) error
	Shutdown() error
}

var backends = map[string]func() Renderer{}

func RegisterBackend(name string, factory func() Renderer) {
	backends[name] = factory
}

func NewBackend(name string) (Renderer, error) {
	f, ok := backends[name]
	if !ok {
		return nil, fmt.Errorf("nativ: unknown renderer backend %q", name)
	}
	return f(), nil
}

func ListBackends() []string {
	names := make([]string, 0, len(backends))
	for n := range backends {
		names = append(names, n)
	}
	return names
}

var ErrNotSupported = errors.New("nativ: renderer not available")
