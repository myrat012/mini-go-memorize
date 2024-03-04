package app

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	_ "github.com/mattn/go-sqlite3"
	"github.com/myrat012/mini-go-memorize/internal/model"
	"github.com/myrat012/mini-go-memorize/internal/sqlite"
)

type Memorize struct {
	App      fyne.App
	Database *sqlite.Sqlite
}

var settingsIsDark, settingsIsRandom bool
var settingsQuestionNumber int

func NewMemorize(dbPath string) (*Memorize, error) {
	app := app.NewWithID("com.example.memorize")
	db, err := sqlite.NewConnection(dbPath)
	if err != nil {
		return nil, err
	}

	return &Memorize{
		App:      app,
		Database: db,
	}, nil
}

func (m *Memorize) Start() {
	// create word table
	err := m.Database.CreateWordTable()
	if err != nil {
		fmt.Println("Error call CreateWordTable.")
		return
	}

	// create settings table
	err = m.Database.CreateSettingTable()
	if err != nil {
		fmt.Println("Error call CreateSettingTable.")
		return
	}

	// set settings
	settings, err := m.Database.SelectSettingsTable()
	if err != nil {
		fmt.Println("Error set settings.")
		return
	}

	settingsIsRandom = settings.IsRandom
	settingsQuestionNumber = settings.Questions

	if settings.DarkTheme {
		settingsIsDark = true
		m.App.Settings().SetTheme(theme.DarkTheme())
	} else {
		settingsIsDark = false
		m.App.Settings().SetTheme(theme.LightTheme())
	}

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

	// tranlsate word
	translatedLabel := canvas.NewText(translatedWordString, color.White)
	translatedLabel.TextSize = 30

	// number of questions
	numberQ := canvas.NewText(fmt.Sprintf("Number of Questions: %d", settingsQuestionNumber), color.White)
	numberQ.TextSize = 15

	// is randomly
	isRandomly := canvas.NewText(fmt.Sprintf("Randomly questions: %t", settingsIsRandom), color.White)
	isRandomly.TextSize = 15

	mainContainer := container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewCenterLayout(), wordLabel),
		translatedLabel,
		widget.NewButton("Show", func() {}),
		widget.NewButton("Next", func() {}),
		numberQ,
		isRandomly,
	)
	mainWindow.SetContent(mainContainer)

	mMenu := mainMenu(mainWindow, mainContainer, m)

	mainWindow.SetMainMenu(fyne.NewMainMenu(mMenu))

	mainWindow.ShowAndRun()
}

func mainMenu(mainWindow fyne.Window, mainContainer *fyne.Container, m *Memorize) *fyne.Menu {
	// settings container
	questionCount := widget.NewEntry()
	radio := widget.NewRadioGroup([]string{"Dark Theme", "Light Theme"}, func(value string) {
		if value == "Dark Theme" {
			settingsIsDark = true
		} else {
			settingsIsDark = false
		}
	})
	randomly := widget.NewCheck("Want Randomly?", func(b bool) {
		if b {
			settingsIsRandom = true
		} else {
			settingsIsRandom = false
		}
	})

	// set default settings
	questionCount.SetText(fmt.Sprintf("%d", settingsQuestionNumber))
	if settingsIsDark {
		radio.SetSelected("Dark Theme")
	} else {
		radio.SetSelected("Light Theme")
	}
	if settingsIsRandom {
		randomly.SetChecked(true)
	} else {
		randomly.SetChecked(false)
	}

	settingsContainer := container.New(
		layout.NewVBoxLayout(),
		randomly,
		questionCount,
		radio,
		widget.NewButton("Save", func() {
			settingsQuestionNumber, err := strconv.Atoi(questionCount.Text)
			if err != nil {
				fmt.Println("Error can't convert to int")
				return
			}
			settings := &model.Settings{
				IsRandom:  settingsIsRandom,
				Questions: settingsQuestionNumber,
				DarkTheme: settingsIsDark,
			}
			err = m.Database.UpdateSettingsTable(settings)
			if err != nil {
				fmt.Println("Error can't update settings table")
				return
			}
			dialog.NewInformation("Alert", "Please reboot program.", mainWindow).Show()
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
			// add words to database
			newWord := model.Dictinary{
				Word:           originalWord.Text,
				TranslatedWord: translatedWord.Text,
			}
			err := m.Database.InsertWordTable(&newWord)
			if err != nil {
				fmt.Println("Error can't Insert word table")
				return
			}

			originalWord.SetText("")
			translatedWord.SetText("")
		}),
	)

	item1 := fyne.NewMenuItem("Main", func() {
		mainWindow.SetContent(mainContainer)
	})

	item2 := fyne.NewMenuItem("Settings", func() {
		mainWindow.SetContent(settingsContainer)
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
