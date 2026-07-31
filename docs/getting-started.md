# Getting Started

## Prerequisites

- Go 1.21+
- Windows 10+ (v0.1)
- GPU with OpenGL 3.3+ drivers

## Install the demo

```powershell
go install github.com/CephandriusMaxtori/Nativ/cmd/demo@latest
demo -renderer=opengl
```

## Build from source

```powershell
git clone https://github.com/CephandriusMaxtori/Nativ.git
cd Nativ
go build ./...
go run ./cmd/demo/ -renderer=opengl
```

## Run the counter example

```powershell
go run ./examples/counter/ -renderer=opengl
```

A window with `+`, `-`, `Reset`, and `Quit` buttons and a count label.

## Use as a dependency

```powershell
go get github.com/CephandriusMaxtori/Nativ@v0.1.0
```

```go
package main

import (
    "fmt"
    nativ "github.com/CephandriusMaxtori/Nativ"
    "github.com/CephandriusMaxtori/Nativ/widget"
)

func main() {
    nativ.Init()
    win, _ := nativ.NewWindow("Hello", 400, 300, "opengl")

    btn := widget.NewButton("Click me")
    btn.OnClick = func() { fmt.Println("clicked") }

    layout := widget.NewVBoxLayout()
    layout.Add(widget.NewLabel("Hello from Nativ"))
    layout.Add(btn)

    win.SetRoot(layout)
    win.Show()
    win.Run()
    nativ.Quit()
}
```

## Run the tests

```powershell
go test -count=1 ./...
```
