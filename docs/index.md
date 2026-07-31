# Nativ

**Pure Go GUI library — zero Cgo dependencies.**

Nativ provides a hardware-accelerated GUI toolkit for Go on Windows, with a
cross-platform architecture designed for Linux/macOS. It targets OpenGL 3.3+
without requiring any external C libraries, GPU SDKs, or framework bindings.

## Features

- **Zero Cgo** — all system APIs (Win32, OpenGL) called via `syscall` at runtime
- **Hardware-accelerated** — OpenGL 3.3+ renderer via `opengl32.dll`
- **Widget system** — Button, Label, BoxLayout (VBox/HBox)
- **Event system** — mouse events, window resize, close
- **Cross-platform architecture** — Platform interface + build tags

## Quick start

```powershell
go install github.com/CephandriusMaxtori/Nativ/cmd/demo@latest
demo -renderer=opengl
```

Or run the counter example:

```powershell
go run github.com/CephandriusMaxtori/Nativ/examples/counter@latest -renderer=opengl
```

## Status

v0.1 — Windows-only, pre-alpha. Linux/macOS stubs return `ErrUnsupported`.

| Component | Status |
|---|---|
| Win32 window + event loop | Stable |
| OpenGL renderer (opengl32.dll) | Stable |
| Widget system (Button, Label, BoxLayout) | Stable |
| Vulkan renderer | Blocked (no WSI surface bindings) |
| Text rendering | Planned (v0.2) |
| Linux (X11/Wayland) | Stub |
| macOS (Cocoa) | Stub |

## Documentation

- [Getting Started](getting-started.md) — install, build, run
- [Architecture](architecture.md) — package layout, design decisions
- [API Reference](api.md) — widgets, window, events, renderer
