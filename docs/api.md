# API Reference

## Package `nativ`

Top-level package for application lifecycle.

### Functions

```go
func Init() error
```
Initialize the Nativ library. Must be called before creating windows.

```go
func Quit()
```
Clean up Nativ resources.

```go
func IsInitialized() bool
```
Check if Init has been called.

```go
func AvailableRenderers() []string
```
List registered renderer backends.

### Type `Window`

```go
func NewWindow(title string, width, height int, rendererName string) (*Window, error)
```
Create a new window with the given title, dimensions, and renderer backend.

```go
func (w *Window) SetRoot(root widget.Widget)
```
Set the root widget of the window's widget tree.

```go
func (w *Window) Root() widget.Widget
```
Get the current root widget.

```go
func (w *Window) Show()
```
Show the window.

```go
func (w *Window) Run() error
```
Run the window's message pump. Blocks until the window is closed.

```go
func (w *Window) Close()
```
Request the window to close.

---

## Package `widget`

### Type `Widget`

```go
type Widget interface {
    Bounds() (x, y, w, h float32)
    SetBounds(x, y, w, h float32)
    Draw(cmds *[]renderer.DrawCommand)
    HandleEvent(evt event.Event) bool
    Parent() Widget
    SetParent(p Widget)
}
```

### Type `Button`

```go
func NewButton(text string) *Button
```
Create a button with the given text.

```go
type Button struct {
    Text    string
    OnClick func()
    // ...
}
```
- `OnClick` is called when the button is clicked (MouseDown + MouseUp inside bounds).
- Visual states: normal (green), hovered (lighter), pressed (darker).

### Type `Label`

```go
func NewLabel(text string) *Label
```
Create a label with the given text.

```go
type Label struct {
    Text string
    // ...
}
```

### Type `BoxLayout`

```go
func NewVBoxLayout() *BoxLayout
```
Create a vertical box layout (children stacked top-to-bottom).

```go
func NewHBoxLayout() *BoxLayout
```
Create a horizontal box layout (children stacked left-to-right).

```go
func (l *BoxLayout) Add(w Widget)
```
Add a child widget.

```go
func (l *BoxLayout) Remove(w Widget)
```
Remove a child widget.

```go
func (l *BoxLayout) Children() []Widget
```
Get the list of children.

Children are sized evenly within the available space, minus `padding` (10)
and `spacing` (8) between children. Layout recalculates when bounds change
or children are added/removed. Events dispatch to children in reverse order
(topmost child first).

---

## Package `event`

### Type `Event`

```go
type Event struct {
    Type         Type
    X, Y         int
    Width, Height int
}
```

### Types

```go
type Type uint8

const (
    MouseDown  Type = iota // Left button pressed
    MouseUp               // Left button released
    MouseMove             // Mouse moved
    MouseClick            // Click completed
    Resize                // Window resized
    Close                 // Window close requested
)
```

---

## Package `renderer`

### Type `DrawCommand`

```go
type DrawCommand struct {
    X, Y, Width, Height float32
    R, G, B, A          float32
}
```
A colored rectangle to render. Coordinates are in window pixel space,
origin at top-left.

### Type `Renderer`

```go
type Renderer interface {
    Init(platform.Platform) error
    Draw([]DrawCommand) error
    Shutdown() error
}
```

### Backend registry

```go
func RegisterBackend(name string, factory func() Renderer)
func NewBackend(name string) (Renderer, error)
func ListBackends() []string
```

---

## Package `internal/platform`

### Type `Platform`

```go
type Platform interface {
    CreateWindow(title string, width, height int) (hwnd, hdc uintptr, err error)
    Pump(eventHandler func(event.Event), render func()) error
    Destroy()
    Hwnd() uintptr
    Hdc() uintptr
    PostQuit()
}
```
