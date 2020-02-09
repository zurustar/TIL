package main

// 同じフォーマットのシートが複数あるエクセルファイルからデータを読み込んで
// SQLiteにつっこむ想定。各シートのフォーマットも適当に想定した。

// go get github.com/360EntSecGroup-Skylar/excelize する必要あり
// なお go get にはgitが必要。
// gccも必要みたい。

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {

	config := LoadConfigurationFile()

	db, err := sql.Open("sqlite3", config.DatabaseName)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	LoadExcel(config.ExcelFile, db)
	LoadCSV(config.CSVFile, db)

}

type Configuration struct {
	DatabaseName string `json:"databaseName"`
	ExcelFile    string `json:"excelFile"`
	CSVFile      string `json:"csvFile"`
}

func NewConfiguration(filename string) Configuration {
	// とりあえずまるごと読む
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	// 解析
	var config Configuration
	err = json.Unmarshal(raw, &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func LoadConfigurationFile() Configuration {
	// デフォルト
	configuration_file := filepath.Join(`.`, `config.json`)

	// なんか渡されてたら上書き
	if len(os.Args) >= 2 {
		configuration_file = os.Args[0]
	}
	return NewConfiguration(configuration_file)
}

func LoadExcel(filename string, db *sql.DB) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// テーブルの準備
	tx.Exec(`DROP TABLE IF EXISTS ExcelTable;`)
	tx.Exec(`CREATE TABLE ExcelTable(C1 TEXT, C2 TEXT, C3 TEXT, i INT, v REAL);`)

	columns := []string{"B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}

	log.Println("loading " + filename)
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
					log.Fatal(err)
				}
				tx.Exec(`INSERT INTO ExcelTable VALUES(?, ?, ?, ?, ?);`, c1, c2, k, i+4, v)
			}
		}
	}
	tx.Commit()
}

func LoadCSV(filename string, db *sql.DB) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// テーブルの準備
	tx.Exec(`DROP TABLE IF EXISTS CSVTable;`)
	tx.Exec(`CREATE TABLE CSVTable(C1 TEXT, C2 TEXT, C3 INT);`)

	log.Println("loading " + filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	reader.Read() // タイトル行をスキップ
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		tx.Exec(`INSERT INTO CSVTable VALUES(?, ?, ?);`, record[0], record[1], record[2])
	}
	tx.Commit()
}
