//go:build windows

package opengl

import (
	"image"
	"strings"

	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type glyphInfo struct {
	x      int
	width  int
	height int
}

type fontAtlas struct {
	texture uint32
	atlasW  int
	atlasH  int
	glyphs  [95]glyphInfo
	glyphH  int
	baseY   int
	pix     []byte
}

func buildDefaultFontAtlas() *fontAtlas {
	advance := 7
	glyphW := 7
	glyphH := 13
	first := 32
	last := 126
	count := last - first + 1

	atlasW := count * advance
	atlasH := glyphH

	img := image.NewAlpha(image.Rect(0, 0, atlasW, atlasH))
	d := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: basicfont.Face7x13,
		Dot:  fixed.P(0, glyphH-3),
	}

	var sb strings.Builder
	for c := first; c <= last; c++ {
		sb.WriteByte(byte(c))
	}
	d.DrawString(sb.String())

	var atlas fontAtlas
	atlas.atlasW = atlasW
	atlas.atlasH = atlasH
	atlas.glyphH = glyphH
	atlas.pix = img.Pix

	for i := 0; i < count; i++ {
		atlas.glyphs[i] = glyphInfo{
			x:      i * advance,
			width:  glyphW,
			height: glyphH,
		}
	}

	return &atlas
}

func buildSfntAtlas(f *renderer.Font) *fontAtlas {
	parsed := f.FontData()
	if parsed == nil {
		return nil
	}

	first := 32
	last := 126
	count := last - first + 1

	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{
		Size:    f.Size(),
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil
	}

	m := face.Metrics()
	glyphH := (m.Ascent + m.Descent).Ceil()
	if glyphH < 1 {
		glyphH = int(f.Size())
	}
	baselineY := m.Ascent.Ceil()

	d := &font.Drawer{Face: face}
	d.Dot = fixed.P(0, baselineY)

	positions := make([]int, count)
	for i := 0; i < count; i++ {
		startX := d.Dot.X
		d.DrawString(string(rune(first + i)))
		endX := d.Dot.X
		_ = endX
		positions[i] = startX.Floor()
	}

	atlasW := d.Dot.X.Floor()
	if atlasW < 1 {
		atlasW = 1
	}

	img := image.NewAlpha(image.Rect(0, 0, atlasW, glyphH))
	d2 := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(0, baselineY),
	}
	for i := 0; i < count; i++ {
		d2.DrawString(string(rune(first + i)))
	}

	var atlas fontAtlas
	atlas.atlasW = atlasW
	atlas.atlasH = glyphH
	atlas.glyphH = glyphH
	atlas.baseY = baselineY
	atlas.pix = img.Pix

	for i := 0; i < count; i++ {
		endPos := atlasW
		if i+1 < count {
			endPos = positions[i+1]
		}
		w := endPos - positions[i]
		if w < 1 {
			w = 1
		}
		atlas.glyphs[i] = glyphInfo{
			x:      positions[i],
			width:  w,
			height: glyphH,
		}
	}

	return &atlas
}

func (fa *fontAtlas) glyphUV(c byte) (u0, v0, u1, v1 float32) {
	if c < 32 || c > 126 {
		c = 32
	}
	g := fa.glyphs[c-32]
	u0 = float32(g.x) / float32(fa.atlasW)
	v0 = 0
	u1 = float32(g.x+g.width) / float32(fa.atlasW)
	v1 = 1
	return
}

func (fa *fontAtlas) measureText(s string) (width, height float32) {
	if len(s) == 0 {
		return 0, 0
	}
	lastIdx := len(s) - 1
	g := fa.glyphs[s[lastIdx]-32]
	endX := g.x + g.width
	return float32(endX), float32(fa.glyphH)
}

func faGlyphWidth() float32 {
	return 7
}

func faGlyphHeight() float32 {
	return 13
}
