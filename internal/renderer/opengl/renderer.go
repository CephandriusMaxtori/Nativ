//go:build windows

package opengl

import (
	"fmt"
	"math"
	"unsafe"

	"nativ/internal/platform"
	"nativ/internal/renderer"
)

func init() {
	renderer.RegisterBackend("opengl", func() renderer.Renderer {
		return &openGLRenderer{}
	})
}

type openGLRenderer struct {
	ctx        *glContext
	width      int
	height     int
	vao        uint32
	vbo        uint32
	program    uint32
	initialized bool
}

const maxGLQuads = 1024

var (
	vertSrc = "#version 330 core\n" +
		"layout(location = 0) in vec2 aPos;\n" +
		"layout(location = 1) in vec4 aColor;\n" +
		"out vec4 vColor;\n" +
		"void main() {\n" +
		"    gl_Position = vec4(aPos, 0.0, 1.0);\n" +
		"    vColor = aColor;\n" +
		"}\x00"

	fragSrc = "#version 330 core\n" +
		"in vec4 vColor;\n" +
		"out vec4 fragColor;\n" +
		"void main() {\n" +
		"    fragColor = vColor;\n" +
		"}\x00"
)

func (r *openGLRenderer) Init(plat platform.Platform) error {
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

	r.width = 800
	r.height = 600
	r.initialized = true
	return nil
}

func (r *openGLRenderer) initGL() error {
	glEnable(GL_BLEND)
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)

	r.program = compileShaderProgram(vertSrc, fragSrc)
	if r.program == 0 {
		return fmt.Errorf("opengl: failed to compile shader program")
	}

	glGenVertexArrays(1, &r.vao)
	glBindVertexArray(r.vao)

	bufferSize := int64(maxGLQuads * 6 * 24) // maxQuads * 6 verts * 6 floats * 4 bytes
	glGenBuffers(1, &r.vbo)
	glBindBuffer(GL_ARRAY_BUFFER, r.vbo)
	glBufferData(GL_ARRAY_BUFFER, bufferSize, nil, GL_DYNAMIC_DRAW)

	glEnableVertexAttribArray(0)
	glVertexAttribPointer(0, 2, GL_FLOAT, false, 24, 0)

	glEnableVertexAttribArray(1)
	glVertexAttribPointer(1, 4, GL_FLOAT, false, 24, 8)

	glBindVertexArray(0)
	return nil
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

	if len(cmds) > maxGLQuads {
		cmds = cmds[:maxGLQuads]
	}

	glUseProgram(r.program)
	glBindVertexArray(r.vao)
	glBindBuffer(GL_ARRAY_BUFFER, r.vbo)

	// Build vertex data
	vertexCount := len(cmds) * 6
	data := make([]float32, 0, vertexCount*6)

	w := float32(width)
	h := float32(height)
	halfW := w * 0.5
	halfH := h * 0.5

	for _, cmd := range cmds {
		x1 := cmd.X
		y1 := cmd.Y
		x2 := cmd.X + cmd.Width
		y2 := cmd.Y + cmd.Height

		nx1 := (x1 - halfW) / halfW
		ny1 := -(y1 - halfH) / halfH
		nx2 := (x2 - halfW) / halfW
		ny2 := -(y2 - halfH) / halfH

		// 6 vertices = 2 triangles forming a quad
		verts := [6]struct{ px, py, r, g, b, a float32 }{
			{nx1, ny1, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx2, ny1, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx2, ny2, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx1, ny1, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx2, ny2, cmd.R, cmd.G, cmd.B, cmd.A},
			{nx1, ny2, cmd.R, cmd.G, cmd.B, cmd.A},
		}
		for _, v := range verts {
			data = append(data, v.px, v.py, v.r, v.g, v.b, v.a)
		}
	}

	// Upload
	bufSize := int64(len(data)) * 4
	if len(data) > 0 {
		glBufferData(GL_ARRAY_BUFFER, bufSize, unsafe.Pointer(&data[0]), GL_DYNAMIC_DRAW)
	}

	glDrawArrays(GL_TRIANGLES, 0, int32(len(data)/6))

	glBindVertexArray(0)
	glUseProgram(0)

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
