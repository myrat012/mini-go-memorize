package model

const SettingsTableName = "tbl_settings"
const WordsTableName = "tbl_words"

type Settings struct {
	IsRandom  bool `db:"is_random"`
	Questions int  `db:"number_questions"`
	DarkTheme bool `db:"dark_theme"`
}

type Dictinary struct {
	Word           string `db:"word"`
	TranslatedWord string `db:"translate_word"`
}
