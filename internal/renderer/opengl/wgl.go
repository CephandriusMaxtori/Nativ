//go:build windows

package opengl

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	libOpengl32 = windows.NewLazySystemDLL("opengl32.dll")
	libGdi32    = windows.NewLazySystemDLL("gdi32.dll")
	libUser32   = windows.NewLazySystemDLL("user32.dll")

	procWglCreateContext  = libOpengl32.NewProc("wglCreateContext")
	procWglMakeCurrent    = libOpengl32.NewProc("wglMakeCurrent")
	procWglDeleteContext  = libOpengl32.NewProc("wglDeleteContext")
	procWglGetProcAddress = libOpengl32.NewProc("wglGetProcAddress")

	procGetDC            = libUser32.NewProc("GetDC")
	procReleaseDC        = libUser32.NewProc("ReleaseDC")
	procSetPixelFormat   = libGdi32.NewProc("SetPixelFormat")
	procChoosePixelFormat = libGdi32.NewProc("ChoosePixelFormat")
	procSwapBuffers      = libGdi32.NewProc("SwapBuffers")
)

const (
	pfdTypeRGBA      = 0
	pfdMainPlane     = 0
	pfdDoubleBuffer  = 1
	pfdDrawToWindow  = 4
	pfdSupportOpenGL = 32
)

type pixelFormatDescriptor struct {
	nSize            uint16
	nVersion         uint16
	dwFlags          uint32
	iPixelType       uint8
	cColorBits       uint8
	cRedBits         uint8
	cRedShift        uint8
	cGreenBits       uint8
	cGreenShift      uint8
	cBlueBits        uint8
	cBlueShift       uint8
	cAlphaBits       uint8
	cAlphaShift      uint8
	cAccumBits       uint8
	cAccumRedBits    uint8
	cAccumGreenBits  uint8
	cAccumBlueBits   uint8
	cAccumAlphaBits  uint8
	cDepthBits       uint8
	cStencilBits     uint8
	cAuxBuffers      uint8
	iLayerType       uint8
	bReserved        uint8
	dwLayerMask      uint32
	dwVisibleMask    uint32
	dwDamageMask     uint32
}

type glContext struct {
	hwnd uintptr
	hdc  uintptr
	hrc  uintptr
}

func createGLContext(hwnd uintptr) (*glContext, error) {
	runtime.LockOSThread()

	hdc, _, err := procGetDC.Call(hwnd)
	if hdc == 0 {
		return nil, fmt.Errorf("opengl: failed to get DC: %w", err)
	}

	pfd := pixelFormatDescriptor{
		nSize:       uint16(unsafe.Sizeof(pixelFormatDescriptor{})),
		nVersion:    1,
		dwFlags:     pfdDrawToWindow | pfdSupportOpenGL | pfdDoubleBuffer,
		iPixelType:  pfdTypeRGBA,
		cColorBits:  32,
		cDepthBits:  24,
		cStencilBits: 8,
		iLayerType:  pfdMainPlane,
	}

	pixelFormat, _, err := procChoosePixelFormat.Call(hdc, uintptr(unsafe.Pointer(&pfd)))
	if pixelFormat == 0 {
		procReleaseDC.Call(hwnd, hdc)
		return nil, fmt.Errorf("opengl: choose pixel format failed: %w", err)
	}

	ret, _, err := procSetPixelFormat.Call(hdc, pixelFormat, uintptr(unsafe.Pointer(&pfd)))
	if ret == 0 {
		procReleaseDC.Call(hwnd, hdc)
		return nil, fmt.Errorf("opengl: set pixel format failed: %w", err)
	}

	hrc, _, err := procWglCreateContext.Call(hdc)
	if hrc == 0 {
		procReleaseDC.Call(hwnd, hdc)
		return nil, fmt.Errorf("opengl: wglCreateContext failed: %w", err)
	}

	made, _, err := procWglMakeCurrent.Call(hdc, hrc)
	if made == 0 {
		procWglDeleteContext.Call(hrc)
		procReleaseDC.Call(hwnd, hdc)
		return nil, fmt.Errorf("opengl: wglMakeCurrent failed: %w", err)
	}

	return &glContext{
		hwnd: hwnd,
		hdc:  hdc,
		hrc:  hrc,
	}, nil
}

func (ctx *glContext) makeCurrent() error {
	made, _, err := procWglMakeCurrent.Call(ctx.hdc, ctx.hrc)
	if made == 0 {
		return fmt.Errorf("opengl: wglMakeCurrent failed: %w", err)
	}
	return nil
}

func (ctx *glContext) swapBuffers() {
	procSwapBuffers.Call(ctx.hdc)
}

func (ctx *glContext) destroy() {
	if ctx.hrc != 0 {
		procWglMakeCurrent.Call(0, 0)
		procWglDeleteContext.Call(ctx.hrc)
	}
	if ctx.hdc != 0 {
		procReleaseDC.Call(ctx.hwnd, ctx.hdc)
	}
}

func wglGetProcAddress(name string) uintptr {
	b, err := syscall.BytePtrFromString(name)
	if err != nil || b == nil {
		return 0
	}
	r, _, _ := procWglGetProcAddress.Call(uintptr(unsafe.Pointer(b)))
	if r != 0 {
		return r
	}
	p := libOpengl32.NewProc(name)
	if p != nil && p.Find() == nil {
		return p.Addr()
	}
	return 0
}
