package app

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Memorize struct {
	App fyne.App
	Log []string
}

func NewMemorize() *Memorize {
	app := app.NewWithID("com.example.memorize")
	return &Memorize{
		App: app,
		Log: []string{},
	}
}

func (m *Memorize) Start() {
	mainWindow := m.App.NewWindow("Memorize")
	mainWindow.SetMaster()
	mainWindow.CenterOnScreen()
	mainWindow.Resize(fyne.NewSize(400, 400))

	// body
	var wordString, translatedWordString string
	wordString = "Word"
	translatedWordString = "Translate"

	// word
	wordLabel := canvas.NewText(wordString, color.White)
	wordLabel.TextStyle.Bold = true
	wordLabel.TextSize = 50

	translatedLabel := canvas.NewText(translatedWordString, color.White)
	translatedLabel.TextSize = 30

	mainContainer := container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewCenterLayout(), wordLabel),
		translatedLabel,
		widget.NewButton("Show", func() {}),
		widget.NewButton("Next", func() {}),
	)
	mainWindow.SetContent(mainContainer)

	mMenu := mainMenu(mainWindow)
	mainWindow.SetMainMenu(fyne.NewMainMenu(mMenu))

	mainWindow.ShowAndRun()
}

func mainMenu(mainWindow fyne.Window) *fyne.Menu {
	item1 := fyne.NewMenuItem("Settings", func() {})
	item2 := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})
	return &fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{item1, item2},
	}
}
