package main

import (
	"flag"
	"fmt"
	"os"

	nativ "github.com/CephandriusMaxtori/Nativ"
	"github.com/CephandriusMaxtori/Nativ/widget"
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

	// Load a custom TrueType font from the Windows font directory.
	customFont, err := nativ.LoadFont(`C:\Windows\Fonts\segoeui.ttf`, 16)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nativ: custom font load failed (using default): %v\n", err)
	}

	bigFont, err := nativ.LoadFont(`C:\Windows\Fonts\segoeui.ttf`, 24)
	if err != nil {
		fmt.Fprintf(os.Stderr, "nativ: big font load failed (using default): %v\n", err)
	}

	defaultLabel := widget.NewLabel("Default bitmap font (basicfont 7x13)")
	segoeLabel := widget.NewLabel("Segoe UI 16px via nativ.LoadFont")
	bigLabel := widget.NewLabel("Segoe UI 24px")
	if customFont != nil {
		segoeLabel.Font = customFont
	}
	if bigFont != nil {
		bigLabel.Font = bigFont
	}

	btn := widget.NewButton("Click Me")
	btn.OnClick = func() {
		fmt.Println("Button clicked!")
	}

	btn2 := widget.NewButton("Segoe UI Button")
	if customFont != nil {
		btn2.Font = customFont
	}

	layout := widget.NewVBoxLayout()
	layout.Add(defaultLabel)
	layout.Add(segoeLabel)
	layout.Add(bigLabel)
	layout.Add(btn)
	layout.Add(btn2)

	win.SetRoot(layout)
	win.Show()

	fmt.Println("Nativ demo running. Close the window to exit.")
	if err := win.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "nativ: run error: %v\n", err)
		os.Exit(1)
	}
}
