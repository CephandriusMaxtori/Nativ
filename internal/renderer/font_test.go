package renderer

import "testing"

func TestNewFontFromDataInvalid(t *testing.T) {
	_, err := NewFontFromData([]byte("not a font"), 14)
	if err == nil {
		t.Fatal("expected error for invalid font data")
	}
}

func TestFontSizeDefault(t *testing.T) {
	_, err := NewFontFromData([]byte("not a font"), 0)
	if err == nil {
		t.Fatal("expected error for invalid font data")
	}
}

func TestFontSizeAccessor(t *testing.T) {
	f := &Font{size: 22}
	if f.Size() != 22 {
		t.Fatalf("expected size 22, got %v", f.Size())
	}
	var nilFont *Font
	if nilFont.Size() != 14 {
		t.Fatalf("expected default size 14, got %v", nilFont.Size())
	}
}

func TestFontDataNil(t *testing.T) {
	var f *Font
	if f.FontData() != nil {
		t.Fatal("expected nil font data for nil font")
	}
}
