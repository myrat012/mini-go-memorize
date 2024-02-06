package app

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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

	mMenu := mainMenu(mainWindow, mainContainer, m)
	mainWindow.SetMainMenu(fyne.NewMainMenu(mMenu))

	mainWindow.ShowAndRun()
}

func mainMenu(mainWindow fyne.Window, mainContainer *fyne.Container, m *Memorize) *fyne.Menu {

	// settings container
	questionCount := widget.NewEntry()
	questionCount.SetPlaceHolder("How much question you want?")
	radio := widget.NewRadioGroup([]string{"Dark Theme", "Light Theme"}, func(value string) {
		if value == "Dark Theme" {
			m.App.Settings().SetTheme(theme.DarkTheme())
		} else {
			m.App.Settings().SetTheme(theme.LightTheme())
		}
	})

	settingsContainer := container.New(
		layout.NewVBoxLayout(),
		widget.NewCheck("Want Randomly?", func(b bool) {
			if b {
				fmt.Println("Checked")
			} else {
				fmt.Println("Uncheked")
			}
		}),
		questionCount,
		radio,
		widget.NewButton("Save", func() {

		}),
	)

	// add word container
	originalWord := widget.NewEntry()
	originalWord.SetPlaceHolder("Word")

	translatedWord := widget.NewEntry()
	translatedWord.SetPlaceHolder("Translate")

	addWordContainer := container.New(
		layout.NewVBoxLayout(),
		originalWord,
		translatedWord,
		widget.NewButton("Save", func() {
			originalWord.SetText("")
			translatedWord.SetText("")
		}),
	)

	item1 := fyne.NewMenuItem("Settings", func() {
		mainWindow.SetContent(settingsContainer)
	})
	item2 := fyne.NewMenuItem("Back", func() {
		mainWindow.SetContent(mainContainer)
	})

	item3 := fyne.NewMenuItem("Add-Words", func() {
		mainWindow.SetContent(addWordContainer)
	})

	item4 := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})
	return &fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{item1, item2, item3, item4},
	}
}
