package model

const SettingsTableName = "tbl_settings"

type Settings struct {
	IsRandom  bool `db:"is_random"`
	Questions int  `db:"number_questions"`
	DarkTheme bool `db:"dark_theme"`
}

type Dictinary struct {
	word           string
	translatedWord string
}
