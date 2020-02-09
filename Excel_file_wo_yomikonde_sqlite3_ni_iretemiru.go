package main

// 同じフォーマットのシートが複数あるエクセルファイルからデータを読み込んで
// SQLiteにつっこむ想定。各シートのフォーマットも適当に想定した。

// go get github.com/360EntSecGroup-Skylar/excelize する必要あり
// なお go get にはgitが必要。
// gccも必要みたい。

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
)

func main() {

	filename := `testdata.xlsx`
	dbname := `testdata.db`

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.Begin();
	if err != nil {
		log.Fatal(err)
	}

	// テーブルの準備
	tx.Exec(`DROP TABLE IF EXISTS TestTable;`)
	tx.Exec(`CREATE TABLE TestTable(C1 TEXT, C2 TEXT, C3 TEXT, i INT, v TEXT);`)

	columns := []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}
	f, err := excelize.OpenFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	for _, sheetName := range f.GetSheetMap() {
		c1, err := f.GetCellValue(sheetName, "A2")
		if err != nil {
			log.Fatal(err)
		}
		c2, err := f.GetCellValue(sheetName, "D2")
		if err != nil {
			log.Fatal(err)
		}
		for r := 5; r <= 9; r += 2 {
			k, err := f.GetCellValue(sheetName, fmt.Sprintf("B%d", r))
			if err != nil {
				log.Fatal(err)
			}
		for i, c := range columns {
				axis := fmt.Sprintf("%s%d", c, r)
				v, err := f.GetCellValue(sheetName, axis)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}
				tx.Exec(`INSERT INTO TestTable VALUES(?, ?, ?, ?, ?);`, c1, c2, k, i+4, v)
			}
		}
	}
	tx.Commit()
	db.Close()
}
