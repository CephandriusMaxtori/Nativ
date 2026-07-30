package main

import (
	"flag"
	"fmt"
	"os"

	nativ "nativ"
	"nativ/widget"
)

func main() {
	rendererFlag := flag.String("renderer", "opengl", "renderer backend (vulkan: pending WSI bindings)")
	flag.Parse()

	if err := nativ.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "nativ: init failed: %v\n", err)
		os.Exit(1)
	}
	defer nativ.Quit()

	win, err := nativ.NewWindow("Nativ Demo", 800, 600, *rendererFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nativ: failed to create window: %v\n", err)
		os.Exit(1)
	}

	btn := widget.NewButton("Click Me")
	btn.OnClick = func() {
		fmt.Println("Button clicked!")
	}

	label := widget.NewLabel("Hello from Nativ v0.1")

	layout := widget.NewVBoxLayout()
	layout.Add(label)
	layout.Add(btn)

	win.SetRoot(layout)
	win.Show()

	fmt.Println("Nativ demo running. Close the window to exit.")
	if err := win.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "nativ: run error: %v\n", err)
		os.Exit(1)
	}
}
