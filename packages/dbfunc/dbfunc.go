package dbfunc

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//Create table if not exist
func createHistTable() error {
	db, err := sql.Open("sqlite3", "./filehistory.db")
	if err != nil {
		return err
	}
	statement := `
	CREATE TABLE hist(
		origin_path TEXT,
		cur_path TEXT PRIMARY KEY
	)
	`
	_, err = db.Exec(statement)
	return err
}

// Initialize
func InitializeDB() {
	if _, err := os.Stat("filehistory.db"); os.IsNotExist(err) {
		_, err := os.Create("filehistory.db")
		if err != nil {
			log.Fatal("建立資料庫錯誤", err)
		}
		err = createHistTable()
		if err != nil {
			log.Fatal("資料庫初始化錯誤: ", err)
		}
	}
}

// Use to record when a file move to trash
func InsertHistory(origin_path, cur_path string) error {
	db, err := sql.Open("sqlite3", "./filehistory.db")
	if err != nil {
		return err
	}
	stmt, err := db.Prepare("REPLACE INTO hist VALUES(?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(origin_path, cur_path)
	if err != nil {
		return err
	}
	log.Printf("Succesfully insert: %s -> %s", origin_path, cur_path)
	return nil
}

// Return original path (before delete) of cur_path
func FindHistory(cur_path string) (string, error) {
	db, err := sql.Open("sqlite3", "./filehistory.db")
	if err != nil {
		return "", err
	}
	row := db.QueryRow("SELECT origin_path FROM hist WHERE cur_path=?", cur_path)
	var origin_path string
	err = row.Scan(&origin_path)
	if err != nil {
		return "", err
	}
	return origin_path, nil
}
