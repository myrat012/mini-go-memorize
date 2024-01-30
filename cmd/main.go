package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {

	mainWord := widget.NewLabelWithStyle("WORD", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	transletedWord := widget.NewLabel("TRANSLATE-WORD")

	app := app.New()
	window := app.NewWindow("Memoraize")
	window.Resize(fyne.NewSize(400, 400))

	btnShow := widget.NewButton("Show-translate", func() {
		transletedWord.SetText("translate")
	})
	btnNext := widget.NewButton("Next-question", func() {
		fmt.Println("question")
	})

	window.SetContent(
		container.NewVBox(
			centeredLabelContainer(mainWord),
			centeredLabelContainer(transletedWord),
			btnShow,
			btnNext,
		),
	)

	window.ShowAndRun()
}

func centeredLabelContainer(label *widget.Label) *fyne.Container {
	return container.New(layout.NewHBoxLayout(), layout.NewSpacer(), label, layout.NewSpacer())
}
