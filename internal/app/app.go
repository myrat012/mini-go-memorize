package app

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	_ "github.com/mattn/go-sqlite3"
	"github.com/myrat012/mini-go-memorize/internal/sqlite"
)

type Memorize struct {
	App      fyne.App
	Database *sqlite.Sqlite
	Screens  ContainerApp
}

var settingsIsDark, settingsIsRandom bool
var settingsQuestionNumber int

func NewMemorize(dbPath string) (*Memorize, error) {
	app := app.NewWithID("com.example.memorize")
	db, err := sqlite.NewConnection(dbPath)
	if err != nil {
		return nil, err
	}

	initContainers := CreateContainersInit(*db)

	return &Memorize{
		App:      app,
		Database: db,
		Screens:  *initContainers,
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
	mainContainer, err := m.Screens.MainContainer()
	if err != nil {
		panic("MainContainer not created")
	}

	mainWindow.SetContent(mainContainer)

	mMenu := mainMenu(mainWindow, mainContainer, m)

	mainWindow.SetMainMenu(fyne.NewMainMenu(mMenu))

	mainWindow.ShowAndRun()
}

// change ifff delete main menu
func mainMenu(mainWindow fyne.Window, mainContainer *fyne.Container, m *Memorize) *fyne.Menu {

	// settings container
	settingsContainer, err := m.Screens.SettingsContainer(mainWindow)
	if err != nil {
		panic("SettingsContainer not created")
	}

	// add-word screen
	addWordContainer, err := m.Screens.AddNewWordContainer()
	if err != nil {
		panic("AddNewWordContainer not created")
	}

	// list word container
	listWordContainer, err := m.Screens.ListWordContainer()
	if err != nil {
		panic("ListWordContainer not created")
	}

	item1 := fyne.NewMenuItem("Main", func() {
		mainWindow.SetContent(mainContainer)
	})

	item2 := fyne.NewMenuItem("Settings", func() {
		mainWindow.SetContent(settingsContainer)
	})

	item3 := fyne.NewMenuItem("Add-Words", func() {
		mainWindow.SetContent(addWordContainer)
	})

	item4 := fyne.NewMenuItem("List-Words", func() {
		mainWindow.SetContent(listWordContainer)
	})

	item5 := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})
	return &fyne.Menu{
		Label: "File",
		Items: []*fyne.MenuItem{item1, item2, item3, item4, item5},
	}
}
