package nativ

import "nativ/internal/renderer"

var initialized bool

func Init() error {
	initialized = true
	return nil
}

func Quit() {
	initialized = false
}

func IsInitialized() bool { return initialized }

func AvailableRenderers() []string {
	return renderer.ListBackends()
}
