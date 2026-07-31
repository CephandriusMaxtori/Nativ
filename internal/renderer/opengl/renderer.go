//go:build windows

package opengl

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/CephandriusMaxtori/Nativ/internal/platform"
	"github.com/CephandriusMaxtori/Nativ/internal/renderer"
)

func init() {
	renderer.RegisterBackend("opengl", func() renderer.Renderer {
		return &openGLRenderer{}
	})
}

type openGLRenderer struct {
	ctx         *glContext
	width       int
	height      int
	vao         uint32
	vbo         uint32
	program     uint32
	whiteTex    uint32
	fontAtlas   *fontAtlas
	fontCache   map[*renderer.Font]*cachedFont
	texUniform  int32
	initialized bool
}

type cachedFont struct {
	atlas   *fontAtlas
	texture uint32
}

type textBatch struct {
	cf    *cachedFont
	start int
	end   int
}

const (
	maxGLQuads  = 1024
	maxGLChars  = 4096
	stride      = 32
	vertFloats  = 8
)

var (
	vertSrc = "#version 330 core\n" +
		"layout(location = 0) in vec2 aPos;\n" +
		"layout(location = 1) in vec2 aUV;\n" +
		"layout(location = 2) in vec4 aColor;\n" +
		"out vec2 vUV;\n" +
		"out vec4 vColor;\n" +
		"void main() {\n" +
		"    gl_Position = vec4(aPos, 0.0, 1.0);\n" +
		"    vUV = aUV;\n" +
		"    vColor = aColor;\n" +
		"}\x00"

	fragSrc = "#version 330 core\n" +
		"uniform sampler2D uTexture;\n" +
		"in vec2 vUV;\n" +
		"in vec4 vColor;\n" +
		"out vec4 fragColor;\n" +
		"void main() {\n" +
		"    float a = texture(uTexture, vUV).r;\n" +
		"    fragColor = vec4(vColor.rgb, vColor.a * a);\n" +
		"}\x00"
)

func (r *openGLRenderer) Init(plat platform.Platform, width, height int) error {
	var err error
	r.ctx, err = createGLContext(plat.Hwnd())
	if err != nil {
		return err
	}

	loadGLProcs()

	if err := r.initGL(); err != nil {
		r.ctx.destroy()
		return err
	}

	r.width = width
	r.height = height
	r.initialized = true
	return nil
}

func (r *openGLRenderer) SetSize(width, height int) {
	r.width = width
	r.height = height
}

func (r *openGLRenderer) initGL() error {
	glEnable(GL_BLEND)
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)

	r.program = compileShaderProgram(vertSrc, fragSrc)
	if r.program == 0 {
		return fmt.Errorf("opengl: failed to compile shader program")
	}

	texLoc := glGetUniformLocation(r.program, &([]byte("uTexture\x00")[0]))
	if texLoc < 0 {
		return fmt.Errorf("opengl: failed to get uTexture uniform location")
	}
	r.texUniform = texLoc

	glUseProgram(r.program)
	glUniform1i(r.texUniform, 0)
	glUseProgram(0)

	glGenVertexArrays(1, &r.vao)
	glBindVertexArray(r.vao)

	bufferSize := int64((maxGLQuads*6 + maxGLChars*6) * vertFloats * 4)
	glGenBuffers(1, &r.vbo)
	glBindBuffer(GL_ARRAY_BUFFER, r.vbo)
	glBufferData(GL_ARRAY_BUFFER, bufferSize, nil, GL_DYNAMIC_DRAW)

	glEnableVertexAttribArray(0)
	glVertexAttribPointer(0, 2, GL_FLOAT, false, stride, 0)

	glEnableVertexAttribArray(1)
	glVertexAttribPointer(1, 2, GL_FLOAT, false, stride, 8)

	glEnableVertexAttribArray(2)
	glVertexAttribPointer(2, 4, GL_FLOAT, false, stride, 16)

	glBindVertexArray(0)

	r.whiteTex = createWhiteTexture()
	if r.whiteTex == 0 {
		return fmt.Errorf("opengl: failed to create white texture")
	}

	r.fontAtlas = buildDefaultFontAtlas()
	r.fontAtlas.texture = createFontAtlasTexture(r.fontAtlas)
	if r.fontAtlas.texture == 0 {
		return fmt.Errorf("opengl: failed to create font atlas texture")
	}
	r.fontCache = map[*renderer.Font]*cachedFont{
		nil: {atlas: r.fontAtlas, texture: r.fontAtlas.texture},
	}

	return nil
}

func (r *openGLRenderer) cachedFontFor(f *renderer.Font) *cachedFont {
	if cf, ok := r.fontCache[f]; ok {
		return cf
	}
	atlas := buildSfntAtlas(f)
	if atlas == nil {
		return r.fontCache[nil]
	}
	atlas.texture = createFontAtlasTexture(atlas)
	if atlas.texture == 0 {
		return r.fontCache[nil]
	}
	cf := &cachedFont{atlas: atlas, texture: atlas.texture}
	r.fontCache[f] = cf
	return cf
}

func createWhiteTexture() uint32 {
	var tex uint32
	glGenTextures(1, &tex)
	glBindTexture(GL_TEXTURE_2D, tex)
	pixel := byte(255)
	glTexImage2D(GL_TEXTURE_2D, 0, 1, 1, 1, 0, GL_RED, GL_UNSIGNED_BYTE, unsafe.Pointer(&pixel))
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE)
	glBindTexture(GL_TEXTURE_2D, 0)
	return tex
}

func createFontAtlasTexture(fa *fontAtlas) uint32 {
	if len(fa.pix) == 0 {
		return 0
	}
	var tex uint32
	glGenTextures(1, &tex)
	glBindTexture(GL_TEXTURE_2D, tex)
	glPixelStorei(GL_UNPACK_ALIGNMENT, 1)
	glTexImage2D(GL_TEXTURE_2D, 0, 1, int32(fa.atlasW), int32(fa.atlasH), 0, GL_RED, GL_UNSIGNED_BYTE, unsafe.Pointer(&fa.pix[0]))
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_S, GL_CLAMP_TO_EDGE)
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_WRAP_T, GL_CLAMP_TO_EDGE)
	glBindTexture(GL_TEXTURE_2D, 0)
	return tex
}

func compileShaderProgram(vertSrc, fragSrc string) uint32 {
	vert := compileShader(GL_VERTEX_SHADER, vertSrc)
	if vert == 0 {
		return 0
	}
	frag := compileShader(GL_FRAGMENT_SHADER, fragSrc)
	if frag == 0 {
		glDeleteShader(vert)
		return 0
	}

	prog := glCreateProgram()
	glAttachShader(prog, vert)
	glAttachShader(prog, frag)
	glLinkProgram(prog)

	var status int32
	glGetProgramiv(prog, GL_LINK_STATUS, &status)
	if status == 0 {
		var logLen int32
		glGetProgramiv(prog, GL_INFO_LOG_LENGTH, &logLen)
		if logLen > 0 {
			log := make([]byte, logLen)
			glGetProgramInfoLog(prog, logLen, nil, &log[0])
		}
		glDeleteProgram(prog)
		prog = 0
	}

	glDeleteShader(vert)
	glDeleteShader(frag)
	return prog
}

func compileShader(stage uint32, src string) uint32 {
	shader := glCreateShader(stage)
	if shader == 0 {
		return 0
	}

	cstr, err := syscallBytePtrFromString(src)
	if err != nil {
		glDeleteShader(shader)
		return 0
	}

	glShaderSource(shader, 1, &cstr, nil)
	glCompileShader(shader)

	var status int32
	glGetShaderiv(shader, GL_COMPILE_STATUS, &status)
	if status == 0 {
		var logLen int32
		glGetShaderiv(shader, GL_INFO_LOG_LENGTH, &logLen)
		if logLen > 0 {
			log := make([]byte, logLen)
			glGetShaderInfoLog(shader, logLen, nil, &log[0])
		}
		glDeleteShader(shader)
		return 0
	}

	return shader
}

func (r *openGLRenderer) Draw(cmds []renderer.DrawCommand) error {
	if !r.initialized {
		return nil
	}

	width, height := r.width, r.height
	if width == 0 || height == 0 {
		return nil
	}

	if err := r.ctx.makeCurrent(); err != nil {
		return err
	}

	glViewport(0, 0, int32(width), int32(height))
	glClearColor(0.1, 0.1, 0.12, 1.0)
	glClear(GL_COLOR_BUFFER_BIT)

	if len(cmds) == 0 {
		r.ctx.swapBuffers()
		return nil
	}

	// Count solid quads and text characters
	solidCount := 0
	charCount := 0
	for _, cmd := range cmds {
		if cmd.Text != "" {
			charCount += len(cmd.Text)
		} else {
			solidCount++
		}
	}
	if solidCount > maxGLQuads {
		solidCount = maxGLQuads
	}
	if charCount > maxGLChars {
		charCount = maxGLChars
	}
	totalVerts := (solidCount + charCount) * 6
	if totalVerts == 0 {
		r.ctx.swapBuffers()
		return nil
	}

	glUseProgram(r.program)
	glBindVertexArray(r.vao)
	glBindBuffer(GL_ARRAY_BUFFER, r.vbo)
	glActiveTexture(GL_TEXTURE0)

	data := make([]float32, 0, totalVerts*vertFloats)

	w := float32(width)
	h := float32(height)
	halfW := w * 0.5
	halfH := h * 0.5

	emitQuad := func(x, y, qw, qh, u0, v0, u1, v1, rCol, gCol, bCol, aCol float32) {
		nx1 := (x - halfW) / halfW
		ny1 := -(y - halfH) / halfH
		nx2 := (x + qw - halfW) / halfW
		ny2 := -(y + qh - halfH) / halfH
		// tri 1
		data = append(data, nx1, ny1, u0, v0, rCol, gCol, bCol, aCol)
		data = append(data, nx2, ny1, u1, v0, rCol, gCol, bCol, aCol)
		data = append(data, nx2, ny2, u1, v1, rCol, gCol, bCol, aCol)
		// tri 2
		data = append(data, nx1, ny1, u0, v0, rCol, gCol, bCol, aCol)
		data = append(data, nx2, ny2, u1, v1, rCol, gCol, bCol, aCol)
		data = append(data, nx1, ny2, u0, v1, rCol, gCol, bCol, aCol)
	}

	// Pass 1: solid quads with white texture
	solidVertCount := 0
	for _, cmd := range cmds {
		if cmd.Text != "" {
			continue
		}
		if solidCount <= 0 {
			continue
		}
		emitQuad(cmd.X, cmd.Y, cmd.Width, cmd.Height,
			0, 0, 0, 0,
			cmd.R, cmd.G, cmd.B, cmd.A)
		solidCount--
		solidVertCount += 6
	}

	// Pass 2: text quads grouped by font
	batches := make([]textBatch, 0, 4)
	batchIdx := make(map[*cachedFont]int)
	for _, cmd := range cmds {
		if cmd.Text == "" {
			continue
		}
		cf := r.cachedFontFor(cmd.Font)
		bi, ok := batchIdx[cf]
		if !ok {
			bi = len(batches)
			batchIdx[cf] = bi
			batches = append(batches, textBatch{cf: cf, start: len(data) / vertFloats, end: len(data) / vertFloats})
		}
		b := &batches[bi]
		atlas := cf.atlas

		textR, textG, textB := cmd.R, cmd.G, cmd.B
		if textR == 0 && textG == 0 && textB == 0 {
			textR, textG, textB = 1, 1, 1
		}
		gh := float32(atlas.glyphH)
		xPos := cmd.X
		for i := 0; i < len(cmd.Text) && charCount > 0; i++ {
			c := cmd.Text[i]
			if c < 32 || c > 126 {
				c = 32
			}
			g := atlas.glyphs[c-32]
			gw := float32(g.width)
			u0, v0, u1, v1 := atlas.glyphUV(c)
			emitQuad(xPos, cmd.Y, gw, gh,
				u0, v0, u1, v1,
				textR, textG, textB, cmd.A)
			xPos += gw
			charCount--
		}
		b.end = len(data) / vertFloats
	}

	if len(data) == 0 {
		r.ctx.swapBuffers()
		return nil
	}

	glBufferData(GL_ARRAY_BUFFER, int64(len(data))*4, unsafe.Pointer(&data[0]), GL_DYNAMIC_DRAW)

	// Draw solid quads (white texture)
	if solidVertCount > 0 {
		glBindTexture(GL_TEXTURE_2D, r.whiteTex)
		glDrawArrays(GL_TRIANGLES, 0, int32(solidVertCount))
	}

	// Draw text quads per font
	offset := solidVertCount
	for _, b := range batches {
		count := b.end - b.start
		if count <= 0 {
			continue
		}
		glBindTexture(GL_TEXTURE_2D, b.cf.texture)
		glDrawArrays(GL_TRIANGLES, int32(offset), int32(count))
		offset += count
	}

	glBindVertexArray(0)
	glUseProgram(0)
	glBindTexture(GL_TEXTURE_2D, 0)

	r.ctx.swapBuffers()
	return nil
}

func (r *openGLRenderer) Shutdown() error {
	if r.program != 0 {
		glDeleteProgram(r.program)
	}
	if r.vbo != 0 {
		glDeleteBuffers(1, &r.vbo)
	}
	if r.vao != 0 {
		glDeleteVertexArrays(1, &r.vao)
	}
	if r.ctx != nil {
		r.ctx.destroy()
	}
	r.initialized = false
	return nil
}

func syscallBytePtrFromString(s string) (*byte, error) {
	return &([]byte(s)[0]), nil
}

// math.Float32bits is used in other files, make sure it works here
var _ = math.Float32bits
