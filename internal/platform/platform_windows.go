//go:build windows

package platform

import (
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/CephandriusMaxtori/Nativ/internal/event"
	"golang.org/x/sys/windows"
)

var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	gdi32  = windows.NewLazySystemDLL("gdi32.dll")

	procGetModuleHandleW   = user32.NewProc("GetModuleHandleW")
	procRegisterClassExW   = user32.NewProc("RegisterClassExW")
	procCreateWindowExW    = user32.NewProc("CreateWindowExW")
	procDestroyWindow      = user32.NewProc("DestroyWindow")
	procDefWindowProcW     = user32.NewProc("DefWindowProcW")
	procPostQuitMessage    = user32.NewProc("PostQuitMessage")
	procPeekMessageW       = user32.NewProc("PeekMessageW")
	procTranslateMessage   = user32.NewProc("TranslateMessage")
	procDispatchMessageW   = user32.NewProc("DispatchMessageW")
	procShowWindow         = user32.NewProc("ShowWindow")
	procUpdateWindow       = user32.NewProc("UpdateWindow")
	procGetDC              = user32.NewProc("GetDC")
	procReleaseDC          = user32.NewProc("ReleaseDC")
	procGetClientRect      = user32.NewProc("GetClientRect")
	procSetWindowLongPtrW  = user32.NewProc("SetWindowLongPtrW")
	procGetWindowLongPtrW  = user32.NewProc("GetWindowLongPtrW")
	procAdjustWindowRect   = user32.NewProc("AdjustWindowRect")
	procGetMessageW        = user32.NewProc("GetMessageW")
	procWaitMessage        = user32.NewProc("WaitMessage")
	procUnregisterClassW   = user32.NewProc("UnregisterClassW")

	procGetStockObject = gdi32.NewProc("GetStockObject")
)

const (
	gwlpUserData     = ^uintptr(20) // GWLP_USERDATA = -21 as unsigned
	pmRemove         = 1
	wmDestroy        = 2
	wmClose          = 16
	wmSize           = 5
	wmLButtonDown    = 513
	wmLButtonUp      = 514
	wmMouseMove      = 512
	wsOverlappedWindow = 0x00000000 | 0x00C00000 | 0x00800000 | 0x00040000 | 0x00020000 | 0x00010000
	csHRedraw        = 0x0002
	csVRedraw        = 0x0001
	csDblClks        = 0x0008
	swNormal         = 1
)

type rect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

type wndClassEx struct {
	size       uint32
	style      uint32
	wndProc    uintptr
	clsExtra   int32
	wndExtra   int32
	instance   windows.Handle
	icon       windows.Handle
	cursor     windows.Handle
	background windows.Handle
	menuName   *uint16
	className  *uint16
	smallIcon  windows.Handle
}

type msg struct {
	hwnd    windows.Handle
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	ptX     int32
	ptY     int32
}

var wndprocCallback = syscall.NewCallback(wndProc)
var currentPlatform *win32Platform

var (
	registeredClass uint16
	className16     *uint16
	registerOnce    sync.Once
)

type win32Platform struct {
	hwnd         windows.Handle
	hdc          windows.Handle
	hinstance    windows.Handle
	eventHandler func(event.Event)
	quit         bool
}

func New() (Platform, error) {
	return &win32Platform{}, nil
}

func (p *win32Platform) Hwnd() uintptr { return uintptr(p.hwnd) }
func (p *win32Platform) Hdc() uintptr  { return uintptr(p.hdc) }
func (p *win32Platform) PostQuit()     { p.quit = true }

func ensureWindowClass(hinstance windows.Handle) (uintptr, error) {
	var regErr error
	registerOnce.Do(func() {
		var err error
		className16, err = windows.UTF16PtrFromString("NativWindow")
		if err != nil {
			regErr = err
			return
		}

		wc := wndClassEx{
			size:      uint32(unsafe.Sizeof(wndClassEx{})),
			style:     csHRedraw | csVRedraw | csDblClks,
			wndProc:   wndprocCallback,
			instance:  hinstance,
			cursor:    loadStandardCursor(),
			className: className16,
		}

		r, _, _ := procGetStockObject.Call(uintptr(5)) // WHITE_BRUSH
		if r != 0 {
			wc.background = windows.Handle(r)
		}

		atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))
		if atom == 0 {
			regErr = err
			return
		}
		registeredClass = uint16(atom)
	})
	return uintptr(registeredClass), regErr
}

func (p *win32Platform) CreateWindow(title string, width, height int) (uintptr, uintptr, error) {
	runtime.LockOSThread()

	hinstance, _, err := procGetModuleHandleW.Call(0)
	if hinstance == 0 {
		return 0, 0, err
	}
	p.hinstance = windows.Handle(hinstance)

	if _, err := ensureWindowClass(p.hinstance); err != nil {
		return 0, 0, err
	}

	titleUTF16, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0, 0, err
	}

	var rct rect
	rct.left = 0
	rct.top = 0
	rct.right = int32(width)
	rct.bottom = int32(height)

	procAdjustWindowRect.Call(
		uintptr(unsafe.Pointer(&rct)),
		wsOverlappedWindow,
		0,
	)

	winWidth := rct.right - rct.left
	winHeight := rct.bottom - rct.top

	currentPlatform = p

	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className16)),
		uintptr(unsafe.Pointer(titleUTF16)),
		wsOverlappedWindow,
		0,
		0,
		uintptr(uint32(winWidth)),
		uintptr(uint32(winHeight)),
		0,
		0,
		uintptr(p.hinstance),
		0,
	)
	if hwnd == 0 {
		currentPlatform = nil
		return 0, 0, err
	}
	p.hwnd = windows.Handle(hwnd)

	procSetWindowLongPtrW.Call(
		hwnd,
		gwlpUserData,
		uintptr(unsafe.Pointer(p)),
	)

	dc, _, _ := procGetDC.Call(hwnd)
	p.hdc = windows.Handle(dc)

	return hwnd, dc, nil
}

func (p *win32Platform) Pump(eventHandler func(event.Event), render func()) error {
	p.eventHandler = eventHandler

	procShowWindow.Call(uintptr(p.hwnd), swNormal)
	procUpdateWindow.Call(uintptr(p.hwnd))

	var m msg
	for !p.quit {
		hadMessages := false
		for {
			r, _, _ := procPeekMessageW.Call(
				uintptr(unsafe.Pointer(&m)),
				0,
				0,
				0,
				pmRemove,
			)
			if r == 0 {
				break
			}
			hadMessages = true
			if m.message == wmDestroy {
				p.quit = true
				break
			}
			procTranslateMessage.Call(uintptr(unsafe.Pointer(&m)))
			procDispatchMessageW.Call(uintptr(unsafe.Pointer(&m)))
		}

		render()

		if !p.quit && !hadMessages {
			procWaitMessage.Call()
		}
	}

	return nil
}

func (p *win32Platform) Destroy() {
	if p.hdc != 0 {
		procReleaseDC.Call(uintptr(p.hwnd), uintptr(p.hdc))
		p.hdc = 0
	}
	if p.hwnd != 0 {
		procDestroyWindow.Call(uintptr(p.hwnd))
		p.hwnd = 0
	}
	if currentPlatform == p {
		currentPlatform = nil
	}
}

func unregisterWindowClass(hinstance windows.Handle) {
	if registeredClass != 0 {
		procUnregisterClassW.Call(
			uintptr(unsafe.Pointer(className16)),
			uintptr(hinstance),
		)
		registeredClass = 0
	}
}

func loadStandardCursor() windows.Handle {
	p := user32.NewProc("LoadCursorW")
	r, _, _ := p.Call(0, uintptr(32512)) // IDC_ARROW
	return windows.Handle(r)
}

func LOWORD(lp uintptr) uint16 {
	return uint16(lp & 0xFFFF)
}

func HIWORD(lp uintptr) uint16 {
	return uint16((lp >> 16) & 0xFFFF)
}

func wndProc(hwnd, msg, wparam, lparam uintptr) uintptr {
	p := currentPlatform
	if p == nil {
		r, _, _ := procDefWindowProcW.Call(hwnd, msg, wparam, lparam)
		return r
	}

	switch uint32(msg) {
	case wmClose:
		p.eventHandler(event.Event{Type: event.Close})
		return 0

	case wmDestroy:
		p.quit = true
		return 0

	case wmSize:
		w := int(LOWORD(lparam))
		h := int(HIWORD(lparam))
		p.eventHandler(event.Event{Type: event.Resize, Width: w, Height: h})
		r, _, _ := procDefWindowProcW.Call(hwnd, msg, wparam, lparam)
		return r

	case wmLButtonDown:
		x := int(int16(LOWORD(lparam)))
		y := int(int16(HIWORD(lparam)))
		p.eventHandler(event.Event{Type: event.MouseDown, X: x, Y: y})

	case wmLButtonUp:
		x := int(int16(LOWORD(lparam)))
		y := int(int16(HIWORD(lparam)))
		p.eventHandler(event.Event{Type: event.MouseUp, X: x, Y: y})

	case wmMouseMove:
		x := int(int16(LOWORD(lparam)))
		y := int(int16(HIWORD(lparam)))
		p.eventHandler(event.Event{Type: event.MouseMove, X: x, Y: y})
	}

	r, _, _ := procDefWindowProcW.Call(hwnd, msg, wparam, lparam)
	return r
}
