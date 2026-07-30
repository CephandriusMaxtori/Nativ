package event

type Type uint8

const (
	MouseDown Type = iota
	MouseUp
	MouseMove
	MouseClick
	Resize
	Close
)

type Event struct {
	Type         Type
	X, Y         int
	Width, Height int
}
