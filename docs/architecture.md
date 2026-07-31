# Architecture

## Package layout

```
nativ/
├── app.go                    # Init/Quit lifecycle
├── backends_windows.go       # Blank-imports OpenGL backend
├── window.go                 # Public Window API

├── internal/
│   ├── event/
│   │   └── event.go          # MouseDown/Up/Move, Resize, Close
│   ├── platform/
│   │   ├── platform.go       # Platform interface
│   │   ├── platform_windows.go   # Win32 implementation
│   │   └── stub.go           # Linux/macOS stub
│   └── renderer/
│       ├── render.go         # Renderer interface + backend registry
│       ├── opengl/
│       │   ├── gl.go         # GL function pointer loading
│       │   ├── renderer.go   # Shader-based quad renderer
│       │   └── wgl.go        # WGL context creation
│       └── vulkan/           # Disabled (needs WSI surface bindings)

├── widget/
│   ├── widget.go             # Widget interface + BaseWidget
│   ├── button.go             # Button with hover/press/click
│   ├── label.go              # Static label
│   └── layout.go             # VBox/HBox layout

└── examples/
    └── counter/              # Interactive counter demo
```

## Design decisions

### Zero Cgo

All Windows API and OpenGL calls use `syscall` and `golang.org/x/sys/windows`
to dynamically load DLLs and call functions at runtime. No C compiler, no GPU
SDK, no external DLLs to ship.

Benefits:
- Single-binary deployment
- Trivial cross-compilation
- No C toolchain in CI
- Easy `go build` / `go test`

### Platform abstraction

The `Platform` interface abstracts OS-specific windowing:
- `CreateWindow`, `Pump`, `Destroy`, `PostQuit`
- Build tags select the right implementation
- Non-Windows builds compile but return `ErrUnsupported` at runtime

### Renderer backends

The `Renderer` interface abstracts GPU rendering:
- `Init`, `Draw`, `Shutdown`
- Backends self-register via `init()` functions
- `RegisterBackend` / `NewBackend` / `ListBackends`

### Widget tree

Widgets form a tree with parent/child relationships:
- `Widget` interface: `Bounds`, `SetBounds`, `Draw`, `HandleEvent`
- `BaseWidget` provides default implementations
- `BoxLayout` distributes space evenly among children
- Events dispatch in reverse Z-order (last child wins)

### Message pump

The Win32 message pump uses `PeekMessageW` (non-blocking) to drain the
message queue each frame, then falls back to `WaitMessage` when idle.
This allows continuous rendering for animations while keeping CPU usage
near zero when the window is idle.
