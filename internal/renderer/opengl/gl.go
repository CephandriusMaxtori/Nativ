//go:build windows

package opengl

import (
	"syscall"
	"unsafe"
)

type GLbitfield uint32
type GLboolean uint8
type GLint int32
type GLuint uint32
type GLsizei int32
type GLfloat float32
type GLclampf float32
type GLsizeiptr int64
type GLchar byte
type GLenum uint32

const (
	GL_FALSE                   = 0
	GL_TRUE                    = 1
	GL_COLOR_BUFFER_BIT        = 0x4000
	GL_DEPTH_BUFFER_BIT        = 0x0100
	GL_BLEND                   = 0x0BE2
	GL_SRC_ALPHA               = 0x0302
	GL_ONE_MINUS_SRC_ALPHA     = 0x0303
	GL_TRIANGLES               = 0x0004
	GL_FLOAT                   = 0x1406
	GL_ARRAY_BUFFER            = 0x8892
	GL_STATIC_DRAW             = 0x88E4
	GL_DYNAMIC_DRAW            = 0x88E8
	GL_VERTEX_SHADER           = 0x8B31
	GL_FRAGMENT_SHADER         = 0x8B30
	GL_COMPILE_STATUS          = 0x8B81
	GL_LINK_STATUS             = 0x8B82
	GL_INFO_LOG_LENGTH         = 0x8B84
)

type glProcs struct {
	viewport              uintptr
	clearColor            uintptr
	clear                 uintptr
	enable                uintptr
	disable               uintptr
	blendFunc             uintptr
	genVertexArrays       uintptr
	bindVertexArray       uintptr
	deleteVertexArrays    uintptr
	genBuffers            uintptr
	bindBuffer            uintptr
	bufferData            uintptr
	deleteBuffers         uintptr
	enableVertexAttribArray  uintptr
	disableVertexAttribArray uintptr
	vertexAttribPointer      uintptr
	createShader          uintptr
	shaderSource          uintptr
	compileShader         uintptr
	getShaderiv           uintptr
	getShaderInfoLog      uintptr
	deleteShader          uintptr
	createProgram         uintptr
	attachShader          uintptr
	linkProgram           uintptr
	getProgramiv          uintptr
	getProgramInfoLog     uintptr
	useProgram            uintptr
	deleteProgram         uintptr
	drawArrays            uintptr
	getError              uintptr
	getUniformLocation    uintptr
	uniform4f             uintptr
	uniformMatrix4fv      uintptr
}

var gl glProcs

func loadGLProcs() {
	gl.viewport = getGLProc("glViewport")
	gl.clearColor = getGLProc("glClearColor")
	gl.clear = getGLProc("glClear")
	gl.enable = getGLProc("glEnable")
	gl.disable = getGLProc("glDisable")
	gl.blendFunc = getGLProc("glBlendFunc")
	gl.genVertexArrays = getGLProc("glGenVertexArrays")
	gl.bindVertexArray = getGLProc("glBindVertexArray")
	gl.deleteVertexArrays = getGLProc("glDeleteVertexArrays")
	gl.genBuffers = getGLProc("glGenBuffers")
	gl.bindBuffer = getGLProc("glBindBuffer")
	gl.bufferData = getGLProc("glBufferData")
	gl.deleteBuffers = getGLProc("glDeleteBuffers")
	gl.enableVertexAttribArray = getGLProc("glEnableVertexAttribArray")
	gl.disableVertexAttribArray = getGLProc("glDisableVertexAttribArray")
	gl.vertexAttribPointer = getGLProc("glVertexAttribPointer")
	gl.createShader = getGLProc("glCreateShader")
	gl.shaderSource = getGLProc("glShaderSource")
	gl.compileShader = getGLProc("glCompileShader")
	gl.getShaderiv = getGLProc("glGetShaderiv")
	gl.getShaderInfoLog = getGLProc("glGetShaderInfoLog")
	gl.deleteShader = getGLProc("glDeleteShader")
	gl.createProgram = getGLProc("glCreateProgram")
	gl.attachShader = getGLProc("glAttachShader")
	gl.linkProgram = getGLProc("glLinkProgram")
	gl.getProgramiv = getGLProc("glGetProgramiv")
	gl.getProgramInfoLog = getGLProc("glGetProgramInfoLog")
	gl.useProgram = getGLProc("glUseProgram")
	gl.deleteProgram = getGLProc("glDeleteProgram")
	gl.drawArrays = getGLProc("glDrawArrays")
	gl.getError = getGLProc("glGetError")
	gl.getUniformLocation = getGLProc("glGetUniformLocation")
	gl.uniform4f = getGLProc("glUniform4f")
	gl.uniformMatrix4fv = getGLProc("glUniformMatrix4fv")
}

func getGLProc(name string) uintptr {
	addr := wglGetProcAddress(name)
	if addr == 0 {
		// Try opengl32.dll directly for standard functions
		p := libOpengl32.NewProc(name)
		if p != nil && p.Find() == nil {
			return p.Addr()
		}
	}
	return addr
}

func glViewport(x, y, width, height int32) {
	syscall.Syscall6(gl.viewport, 4, uintptr(x), uintptr(y), uintptr(width), uintptr(height), 0, 0)
}

func glClearColor(r, g, b, a float32) {
	syscall.Syscall6(gl.clearColor, 4,
		uintptr(r), uintptr(g), uintptr(b), uintptr(a), 0, 0)
}

func glClear(mask uint32) {
	syscall.Syscall(gl.clear, 1, uintptr(mask), 0, 0)
}

func glEnable(cap uint32) {
	syscall.Syscall(gl.enable, 1, uintptr(cap), 0, 0)
}

func glDisable(cap uint32) {
	syscall.Syscall(gl.disable, 1, uintptr(cap), 0, 0)
}

func glBlendFunc(sfactor, dfactor uint32) {
	syscall.Syscall(gl.blendFunc, 2, uintptr(sfactor), uintptr(dfactor), 0)
}

func glGenVertexArrays(n int32, arrays *uint32) {
	syscall.Syscall(gl.genVertexArrays, 2, uintptr(n), uintptr(unsafe.Pointer(arrays)), 0)
}

func glBindVertexArray(array uint32) {
	syscall.Syscall(gl.bindVertexArray, 1, uintptr(array), 0, 0)
}

func glDeleteVertexArrays(n int32, arrays *uint32) {
	syscall.Syscall(gl.deleteVertexArrays, 2, uintptr(n), uintptr(unsafe.Pointer(arrays)), 0)
}

func glGenBuffers(n int32, buffers *uint32) {
	syscall.Syscall(gl.genBuffers, 2, uintptr(n), uintptr(unsafe.Pointer(buffers)), 0)
}

func glBindBuffer(target, buffer uint32) {
	syscall.Syscall(gl.bindBuffer, 2, uintptr(target), uintptr(buffer), 0)
}

func glBufferData(target uint32, size int64, data unsafe.Pointer, usage uint32) {
	syscall.Syscall6(gl.bufferData, 4, uintptr(target), uintptr(size), uintptr(data), uintptr(usage), 0, 0)
}

func glDeleteBuffers(n int32, buffers *uint32) {
	syscall.Syscall(gl.deleteBuffers, 2, uintptr(n), uintptr(unsafe.Pointer(buffers)), 0)
}

func glEnableVertexAttribArray(index uint32) {
	syscall.Syscall(gl.enableVertexAttribArray, 1, uintptr(index), 0, 0)
}

func glDisableVertexAttribArray(index uint32) {
	syscall.Syscall(gl.disableVertexAttribArray, 1, uintptr(index), 0, 0)
}

func glVertexAttribPointer(index uint32, size int32, xtype uint32, normalized bool, stride int32, offset uintptr) {
	norm := uint32(0)
	if normalized {
		norm = 1
	}
	syscall.Syscall6(gl.vertexAttribPointer, 6,
		uintptr(index), uintptr(size), uintptr(xtype), uintptr(norm), uintptr(stride), offset)
}

func glCreateShader(xtype uint32) uint32 {
	r, _, _ := syscall.Syscall(gl.createShader, 1, uintptr(xtype), 0, 0)
	return uint32(r)
}

func glShaderSource(shader uint32, count int32, str **byte, length *int32) {
	syscall.Syscall6(gl.shaderSource, 4,
		uintptr(shader), uintptr(count), uintptr(unsafe.Pointer(str)), uintptr(unsafe.Pointer(length)), 0, 0)
}

func glCompileShader(shader uint32) {
	syscall.Syscall(gl.compileShader, 1, uintptr(shader), 0, 0)
}

func glGetShaderiv(shader uint32, pname uint32, params *int32) {
	syscall.Syscall(gl.getShaderiv, 3, uintptr(shader), uintptr(pname), uintptr(unsafe.Pointer(params)))
}

func glGetShaderInfoLog(shader uint32, bufSize int32, length *int32, infoLog *byte) {
	syscall.Syscall6(gl.getShaderInfoLog, 4,
		uintptr(shader), uintptr(bufSize), uintptr(unsafe.Pointer(length)), uintptr(unsafe.Pointer(infoLog)), 0, 0)
}

func glDeleteShader(shader uint32) {
	syscall.Syscall(gl.deleteShader, 1, uintptr(shader), 0, 0)
}

func glCreateProgram() uint32 {
	r, _, _ := syscall.Syscall(gl.createProgram, 0, 0, 0, 0)
	return uint32(r)
}

func glAttachShader(program, shader uint32) {
	syscall.Syscall(gl.attachShader, 2, uintptr(program), uintptr(shader), 0)
}

func glLinkProgram(program uint32) {
	syscall.Syscall(gl.linkProgram, 1, uintptr(program), 0, 0)
}

func glGetProgramiv(program uint32, pname uint32, params *int32) {
	syscall.Syscall(gl.getProgramiv, 3, uintptr(program), uintptr(pname), uintptr(unsafe.Pointer(params)))
}

func glGetProgramInfoLog(program uint32, bufSize int32, length *int32, infoLog *byte) {
	syscall.Syscall6(gl.getProgramInfoLog, 4,
		uintptr(program), uintptr(bufSize), uintptr(unsafe.Pointer(length)), uintptr(unsafe.Pointer(infoLog)), 0, 0)
}

func glUseProgram(program uint32) {
	syscall.Syscall(gl.useProgram, 1, uintptr(program), 0, 0)
}

func glDeleteProgram(program uint32) {
	syscall.Syscall(gl.deleteProgram, 1, uintptr(program), 0, 0)
}

func glDrawArrays(mode uint32, first int32, count int32) {
	syscall.Syscall(gl.drawArrays, 3, uintptr(mode), uintptr(first), uintptr(count))
}

func glGetError() uint32 {
	r, _, _ := syscall.Syscall(gl.getError, 0, 0, 0, 0)
	return uint32(r)
}

func glGetUniformLocation(program uint32, name *byte) int32 {
	r, _, _ := syscall.Syscall(gl.getUniformLocation, 2, uintptr(program), uintptr(unsafe.Pointer(name)), 0)
	return int32(r)
}

func glUniform4f(location int32, v0, v1, v2, v3 float32) {
	syscall.Syscall6(gl.uniform4f, 5,
		uintptr(location), uintptr(v0), uintptr(v1), uintptr(v2), uintptr(v3), 0)
}

func glUniformMatrix4fv(location int32, count int32, transpose bool, value *float32) {
	t := uint32(0)
	if transpose {
		t = 1
	}
	syscall.Syscall6(gl.uniformMatrix4fv, 4,
		uintptr(location), uintptr(count), uintptr(t), uintptr(unsafe.Pointer(value)), 0, 0)
}
