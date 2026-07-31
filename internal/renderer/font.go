package renderer

import (
	"fmt"
	"os"

	"golang.org/x/image/font/sfnt"
)

type Font struct {
	parsed *sfnt.Font
	size   float64
}

func NewFontFromData(data []byte, size float64) (*Font, error) {
	parsed, err := sfnt.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("nativ: invalid font data: %w", err)
	}
	if size <= 0 {
		size = 14
	}
	return &Font{parsed: parsed, size: size}, nil
}

func NewFontFromFile(path string, size float64) (*Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("nativ: %w", err)
	}
	return NewFontFromData(data, size)
}

func (f *Font) Size() float64 {
	if f == nil {
		return 14
	}
	return f.size
}

func (f *Font) FontData() *sfnt.Font {
	if f == nil {
		return nil
	}
	return f.parsed
}
