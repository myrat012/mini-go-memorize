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
	"fyne.io/fyne/v2/widget"
	"github.com/myrat012/mini-go-memorize/internal/model"
	"github.com/myrat012/mini-go-memorize/internal/sqlite"
)

const (
	newWord          = "New Word"
	newWordTranslate = "Word Translate"
	saveBtnName      = "Save"
	emtyString       = ""
)

type ContainerApp struct {
	Database sqlite.Sqlite
}

func CreateContainersInit(db sqlite.Sqlite) *ContainerApp {
	return &ContainerApp{
		Database: db,
	}
}

func (c *ContainerApp) AddNewWordContainer() (*fyne.Container, error) {
	originalWord := widget.NewEntry()
	originalWord.SetPlaceHolder(newWord)

	translatedWord := widget.NewEntry()
	translatedWord.SetPlaceHolder(newWordTranslate)

	box := container.New(
		layout.NewVBoxLayout(),
		originalWord,
		translatedWord,
		widget.NewButton(saveBtnName, func() {

			// add words to database
			newWord := model.Dictinary{
				Word:           originalWord.Text,
				TranslatedWord: translatedWord.Text,
			}
			err := c.Database.InsertWordTable(&newWord)
			if err != nil {
				fmt.Println("Error can't Insert word table")
				return
			}

			originalWord.SetText(emtyString)
			translatedWord.SetText(emtyString)
		}),
	)

	return box, nil
}

func (c *ContainerApp) SettingsContainer(mainWindow fyne.Window) (*fyne.Container, error) {

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
			settings := model.Settings{
				IsRandom:  settingsIsRandom,
				Questions: settingsQuestionNumber,
				DarkTheme: settingsIsDark,
			}
			err = c.Database.UpdateSettingsTable(&settings)
			if err != nil {
				fmt.Println("Error can't update settings table")
				return
			}
			dialog.NewInformation("Alert", "Please reboot program.", mainWindow).Show()
		}),
	)
	return settingsContainer, nil
}

func (c *ContainerApp) MainContainer() (*fyne.Container, error) {
	var wordString, translatedWordString string

	w, err := c.Database.SelectWordTable(settingsQuestionNumber)
	if err != nil {
		fmt.Println("Error SelectWordTable")
		return nil, err
	}
	count := 0

	wordString = ""
	translatedWordString = ""

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
		widget.NewButton("Show", func() {
			translatedLabel.Text = w[count-1]["tword"]
			translatedLabel.Refresh()
		}),
		widget.NewButton("Next", func() {
			if count > len(w) {
				return
			}
			wordLabel.Text = w[count]["word"]
			wordLabel.Refresh()
			count++
			translatedLabel.Text = ""
			translatedLabel.Refresh()
		}),
		numberQ,
		isRandomly,
	)
	return mainContainer, nil
}

func (c *ContainerApp) ListWordContainer() (*fyne.Container, error) {
	idInput := widget.NewEntry()
	idInput.SetPlaceHolder("Id")

	listWord, err := c.Database.ListWordTable()
	if err != nil {
		fmt.Println("Can't get List of words")
		return nil, err
	}

	a := app.New()
	secondScreen := a.NewWindow("table")
	secondScreen.Resize(fyne.NewSize(500, 500))

	table := widget.NewTable(
		func() (rows int, cols int) {
			return len(listWord), 3
		},
		func() fyne.CanvasObject {
			l := widget.NewLabel("Emty")
			return l
		},
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			i := co.(*widget.Label)
			i.SetText(fmt.Sprint(listWord[tci.Row][tci.Col]))
		},
	)

	listWordContainer := container.New(
		layout.NewVBoxLayout(),
		idInput,
		widget.NewButton("Delete row", func() {
			id, err := strconv.Atoi(idInput.Text)
			if err != nil {
				fmt.Println("Error can't convert to int")
				return
			}

			err = c.Database.DeleteWordTable(id)
			if err != nil {
				fmt.Println("Error delete data from word table")
				return
			}
		}),
		widget.NewButton("show table", func() {
			secondScreen.SetContent(table)
			secondScreen.Show()
		}),
	)

	return listWordContainer, nil
}
