//go:build !windows

package platform

import (
	"errors"

	"nativ/internal/event"
)

var ErrUnsupported = errors.New("nativ: windows is the only supported platform in v0.1")

type stubPlatform struct{}

func New() (Platform, error) { return &stubPlatform{}, nil }

func (p *stubPlatform) CreateWindow(title string, width, height int) (uintptr, uintptr, error) {
	return 0, 0, ErrUnsupported
}
func (p *stubPlatform) Pump(func(event.Event), func()) error { return ErrUnsupported }
func (p *stubPlatform) Destroy()                             {}
func (p *stubPlatform) Hwnd() uintptr                        { return 0 }
func (p *stubPlatform) Hdc() uintptr                         { return 0 }
func (p *stubPlatform) PostQuit()                             {}
