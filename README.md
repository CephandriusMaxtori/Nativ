# Nativ

![WIP](https://img.shields.io/badge/status-WIP-yellow)
[![Docs](https://img.shields.io/badge/docs-github.io-blue)](https://cephandriusmaxtori.github.io/Nativ/)

**Pure Go GUI library — zero Cgo dependencies.**

Nativ provides a hardware-accelerated GUI toolkit for Go on Windows, with a cross-platform architecture designed for Linux/macOS. It targets OpenGL 3.3+ without requiring any external C libraries, GPU SDKs, or framework bindings. It includes text rendering with built-in bitmap fonts and custom TTF/OTF fonts at any size.

## Status

**v0.2 — Pre-alpha.** Works on Windows (OpenGL 3.3+). Linux/macOS stubs return `ErrUnsupported`.

| Component | Status |
|---|---|
| Win32 window + event loop | Done |
| OpenGL renderer (opengl32.dll) | Done |
| OpenGL 3.3 shader-based quad rendering | Done |
| Text rendering (bitmap + custom TTF/OTF) | Done |
| Widget system (Button, Label, BoxLayout) | Done |
| Widget hit-testing + event dispatch | Done |
| Entry widget (text input) | In progress |
| Vulkan renderer | Written, blocked on WSI surface bindings |
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
├── app.go / window.go       # Public API (incl. LoadFont/LoadFontFromData)
├── internal/
│   ├── event/               # MouseDown/Up/Move, Close, Resize
│   ├── platform/            # Platform interface + Win32 impl + stubs
│   └── renderer/
│       ├── render.go        # Renderer interface + DrawCommand + registry
│       ├── font.go          # Font type (TTF/OTF via sfnt/opentype)
│       ├── opengl/          # OpenGL 3.3 backend (Windows) + glyph atlases
│       └── vulkan/          # Vulkan backend (disabled, see below)
├── widget/                  # Button, Label, BoxLayout (platform-agnostic)
└── docs/                    # GitHub Pages site
```

All system APIs (Win32, OpenGL) are called via `syscall`/`golang.org/x/sys/windows` at runtime — no `#cgo`, no C compiler, no GPU SDK required.

## Text rendering & custom fonts

Labels and buttons render text through the GPU. A built-in 7x13 bitmap font is used by default; any TTF/OTF font can be loaded and assigned to a widget:

```go
custom, err := nativ.LoadFont("C:\\Windows\\Fonts\\segoeui.ttf", 16)
if err != nil { log.Fatal(err) }

label := widget.NewLabel("Hello")
label.Font = custom
```

See [docs/fonts.html](https://cephandriusmaxtori.github.io/Nativ/fonts.html) for details.

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
