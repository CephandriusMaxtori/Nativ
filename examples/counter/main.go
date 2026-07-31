package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	nativ "github.com/CephandriusMaxtori/Nativ"
	"github.com/CephandriusMaxtori/Nativ/widget"
)

func main() {
	rendererFlag := flag.String("renderer", "opengl", "renderer backend")
	flag.Parse()

	if err := nativ.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "counter: init failed: %v\n", err)
		os.Exit(1)
	}
	defer nativ.Quit()

	win, err := nativ.NewWindow("Counter", 400, 200, *rendererFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "counter: failed to create window: %v\n", err)
		os.Exit(1)
	}

	count := 0
	countLabel := widget.NewLabel("Count: 0")

	decBtn := widget.NewButton("-")
	decBtn.OnClick = func() {
		count--
		countLabel.Text = "Count: " + strconv.Itoa(count)
	}

	incBtn := widget.NewButton("+")
	incBtn.OnClick = func() {
		count++
		countLabel.Text = "Count: " + strconv.Itoa(count)
	}

	resetBtn := widget.NewButton("Reset")
	resetBtn.OnClick = func() {
		count = 0
		countLabel.Text = "Count: 0"
	}

	quitBtn := widget.NewButton("Quit")
	quitBtn.OnClick = func() {
		win.Close()
	}

	buttonRow := widget.NewHBoxLayout()
	buttonRow.Add(decBtn)
	buttonRow.Add(incBtn)

	controlRow := widget.NewHBoxLayout()
	controlRow.Add(resetBtn)
	controlRow.Add(quitBtn)

	root := widget.NewVBoxLayout()
	root.Add(countLabel)
	root.Add(buttonRow)
	root.Add(controlRow)

	win.SetRoot(root)
	win.Show()

	fmt.Println("Counter running. Close the window or click Quit to exit.")
	if err := win.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "counter: run error: %v\n", err)
		os.Exit(1)
	}
}
