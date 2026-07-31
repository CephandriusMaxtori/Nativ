package nativ

import "github.com/CephandriusMaxtori/Nativ/internal/renderer"

type Font = renderer.Font

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

func LoadFont(path string, size float64) (*Font, error) {
	return renderer.NewFontFromFile(path, size)
}

func LoadFontFromData(data []byte, size float64) (*Font, error) {
	return renderer.NewFontFromData(data, size)
}
