package renderer

import (
	"errors"
	"fmt"
	"sync"

	"github.com/CephandriusMaxtori/Nativ/internal/platform"
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

var (
	backends   = map[string]func() Renderer{}
	backendsMu sync.RWMutex
)

func RegisterBackend(name string, factory func() Renderer) {
	backendsMu.Lock()
	backends[name] = factory
	backendsMu.Unlock()
}

func NewBackend(name string) (Renderer, error) {
	backendsMu.RLock()
	f, ok := backends[name]
	backendsMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("nativ: unknown renderer backend %q", name)
	}
	return f(), nil
}

func ListBackends() []string {
	backendsMu.RLock()
	names := make([]string, 0, len(backends))
	for n := range backends {
		names = append(names, n)
	}
	backendsMu.RUnlock()
	return names
}

var ErrNotSupported = errors.New("nativ: renderer not available")
