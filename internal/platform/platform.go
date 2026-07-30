package platform

import "github.com/CephandriusMaxtori/Nativ/internal/event"

type Platform interface {
	CreateWindow(title string, width, height int) (hwnd, hdc uintptr, err error)
	Pump(eventHandler func(event.Event), render func()) error
	Destroy()
	Hwnd() uintptr
	Hdc() uintptr
	PostQuit()
}
