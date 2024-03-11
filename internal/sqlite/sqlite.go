package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/myrat012/mini-go-memorize/internal/model"
)

type Sqlite struct {
	Database *sql.DB
}

func NewConnection(path string) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Database: db,
	}, nil
}

func (s *Sqlite) CreateSettingTable() (err error) {
	_, err = s.Database.Exec("CREATE TABLE IF NOT EXISTS " + model.SettingsTableName + " (id INTEGER PRIMARY KEY, is_random BOOL, number_questions INTEGER, dark_theme BOOL)")
	if err != nil {
		fmt.Printf("error Creating settings table")
		return err
	}

	isEmty, err := isTableEmty(s.Database, model.SettingsTableName)
	if err != nil {
		fmt.Printf("error isTableEmty CreateSettingTable")
		return err
	}
	if isEmty {
		_, err = s.Database.Exec("INSERT INTO "+model.SettingsTableName+" (is_random, number_questions, dark_theme) VALUES (?,?,?)", true, 0, true)
		if err != nil {
			fmt.Printf("error when Inserting data in SettingTable")
			return err
		}
	}

	return nil
}

func (s *Sqlite) UpdateSettingsTable(settings *model.Settings) (err error) {
	query, err := s.Database.Prepare("UPDATE " + model.SettingsTableName + " SET is_random=?, number_questions=?, dark_theme=? WHERE ID=1")
	if err != nil {
		fmt.Printf("error when Update data in SettingTable")
		return err
	}
	defer query.Close()

	_, err = query.Exec(&settings.IsRandom, &settings.Questions, &settings.DarkTheme)
	if err != nil {
		fmt.Printf("error when Exec data in SettingTable")
		return err
	}

	return nil
}

func (s *Sqlite) SelectSettingsTable() (*model.Settings, error) {
	rows, err := s.Database.Query("SELECT is_random, number_questions, dark_theme FROM " + model.SettingsTableName)
	if err != nil {
		fmt.Println("Error SelectSettingsTable")
		return nil, err
	}
	defer rows.Close()

	var settings model.Settings

	for rows.Next() {
		err = rows.Scan(&settings.IsRandom, &settings.Questions, &settings.DarkTheme)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
	}
	return &settings, nil
}

func (s *Sqlite) SelectWordTable(number int) (map[int]map[string]string, error) {
	strQuery := fmt.Sprintf("SELECT word, translate_word FROM %s ORDER BY RANDOM() LIMIT %d;", model.WordsTableName, number)
	rows, err := s.Database.Query(strQuery)
	if err != nil {
		fmt.Println("Error SelectWordTable")
		return nil, err
	}
	defer rows.Close()

	var word model.ListDictinary
	count := 0
	box := make(map[int]map[string]string)

	for rows.Next() {
		err = rows.Scan(&word.Word, &word.TranslatedWord)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		box[count] = map[string]string{"word": word.Word, "tword": word.TranslatedWord}
		count++
	}

	return box, nil
}

func (s *Sqlite) CreateWordTable() (err error) {
	_, err = s.Database.Exec("CREATE TABLE IF NOT EXISTS " + model.WordsTableName + " (id INTEGER PRIMARY KEY, word VARCHAR(150), translate_word VARCHAR(150))")
	if err != nil {
		fmt.Printf("error Creating Word table")
		return err
	}
	return nil
}

func (s *Sqlite) InsertWordTable(w *model.Dictinary) (err error) {
	_, err = s.Database.Exec("INSERT INTO "+model.WordsTableName+" (word, translate_word) VALUES (?, ?)", &w.Word, &w.TranslatedWord)
	if err != nil {
		fmt.Printf("error Insert Word table")
		return err
	}
	return nil
}

func (s *Sqlite) ListWordTable() ([][]string, error) {
	list, err := s.Database.Query("SELECT * FROM " + model.WordsTableName)
	if err != nil {
		fmt.Println("Error ListWordTable")
		return nil, err
	}
	defer list.Close()

	var listWords [][]string

	for list.Next() {
		var dic model.ListDictinary
		err = list.Scan(&dic.Id, &dic.Word, &dic.TranslatedWord)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, err
		}
		listWords = append(listWords, []string{dic.Id, dic.Word, dic.TranslatedWord})

	}
	return listWords, err
}

func (s *Sqlite) DeleteWordTable(id int) error {
	_, err := s.Database.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=%d", model.WordsTableName, id))
	if err != nil {
		return err
	}
	return nil
}

func isTableEmty(db *sql.DB, tableName string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM " + tableName).Scan(&count)
	if err != nil {
		fmt.Println("Error isTableEmty")
		return false, err
	}
	return count == 0, nil
}
