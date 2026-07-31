# Nativ

![WIP](https://img.shields.io/badge/status-WIP-yellow)
[![Docs](https://img.shields.io/badge/docs-github.io-blue)](https://cephandriusmaxtori.github.io/Nativ/)

**Pure Go GUI library — zero Cgo dependencies.**

Nativ provides a hardware-accelerated GUI toolkit for Go on Windows, with a cross-platform architecture designed for Linux/macOS. It targets OpenGL 3.3+ without requiring any external C libraries, GPU SDKs, or framework bindings.

## Status

**v0.1 — Pre-alpha.** Works on Windows (OpenGL 3.3+). Linux/macOS stubs return `ErrUnsupported`.

| Component | Status |
|---|---|
| Win32 window + event loop | Done |
| OpenGL renderer (opengl32.dll) | Done |
| OpenGL 3.3 shader-based quad rendering | Done |
| Widget system (Button, Label, BoxLayout) | Done |
| Widget hit-testing + event dispatch | Done |
| Vulkan renderer | Written, blocked on WSI surface bindings |
| Text rendering | Planned (v0.2) |
| Linux (X11/Wayland) | Stub |
| macOS (Cocoa) | Stub |

## Getting started

```powershell
go install nativ@latest
# or clone and run the demo:
git clone https://github.com/CephandriusMaxtori/Nativ.git
cd Nativ
go run ./cmd/demo/ -renderer=opengl
```

## Architecture

```
nativ/
├── app.go / window.go       # Public API
├── internal/
│   ├── event/               # MouseDown/Up/Move, Close, Resize
│   ├── platform/            # Platform interface + Win32 impl + stubs
│   └── renderer/
│       ├── render.go        # Renderer interface + backend registry
│       ├── opengl/          # OpenGL 3.3 backend (Windows)
│       └── vulkan/          # Vulkan backend (disabled, see below)
└── widget/                  # Button, Label, BoxLayout (platform-agnostic)
```

All system APIs (Win32, OpenGL) are called via `syscall`/`golang.org/x/sys/windows` at runtime — no `#cgo`, no C compiler, no GPU SDK required.

## Why "no Cgo"?

- Single-binary deployment — no DLLs to ship
- Cross-compilation works trivially
- No C toolchain in CI
- Easy `go build` / `go test` on any platform

## Vulkan

The Vulkan renderer is written but **disabled** (`//go:build ignore`). The `christerso/vulkan-go` fork compiles cleanly but doesn't generate `VK_KHR_win32_surface` extension bindings, making on-screen rendering impossible. It will be re-enabled when a `vkCreateWin32SurfaceKHR` binding is available.

## Build & test

```powershell
go build ./...
go vet ./...
go test -count=1 ./...
```
