package sql

import (
	"Imdb/imdblog"
	"database/sql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"strconv"

	imdb "github.com/eefret/go-imdb"
)
var FileObj = "sqlite"

type TblImdb	 struct {
	gorm.Model
	Title    		string
	ReleasedYear 	string
	Rating      	int
	Genres      	string

}
const dbstr = "imdb.db"
func SqliteMigrate() {

	db, err := gorm.Open(sqlite.Open(dbstr), &gorm.Config{})
	imdblog.WriteFile(FileObj, "creating database "+dbstr, os.Args[0])
	if err != nil {
		imdblog.WriteFile(FileObj, "Migrate "+err.Error(), os.Args[0])
		imdblog.WriteFile(FileObj, "failed to connect database", os.Args[0])
		panic("failed to connect database")
	}
	imdblog.WriteFile(FileObj, "created database "+dbstr, os.Args[0])

	err = db.AutoMigrate(&TblImdb{})
	if err != nil {
		imdblog.WriteFile(FileObj, "failed to create  table", os.Args[0])
		panic("failed to create table")
	}


}

func SqliteCreateData(table interface{}) TblImdb {
	db, err := gorm.Open(sqlite.Open(dbstr), &gorm.Config{})
	if err != nil {
		imdblog.WriteFile(FileObj, "failed to connect database", os.Args[0])
		panic("failed to connect database")
	}
	imdblog.WriteFile(FileObj, "CreateData__ connected database", os.Args[0])
	db.Create(table)
	var tablMap TblImdb
	db.Last(&tablMap)
	return tablMap
}

func SqliteUpdateData(tableData interface{},tableCondition interface{})bool{
	db, err := gorm.Open(sqlite.Open(dbstr), &gorm.Config{})
	if err != nil {
		imdblog.WriteFile(FileObj, "failed to connect database", os.Args[0])
		panic("failed to connect database")
	}
	imdblog.WriteFile(FileObj, "UpdateData connected database", os.Args[0])
	constring := tableCondition.(*TblImdb)
	db.Model(TblImdb{}).Where("title = ?", constring.Title).Updates(tableData)
	if db.Error != nil{
		return false
	}
	return true
}
func Search(condition ,conditionString, conditionValue string)  ([]*TblImdb, bool) {
	var movies TblImdb

	db, err := gorm.Open(sqlite.Open(dbstr), &gorm.Config{})
	if err != nil {
		imdblog.WriteFile(FileObj, "failed to connect database", os.Args[0])
		panic("failed to connect database")
	}
	imdblog.WriteFile(FileObj, "TblImdb connected database", os.Args[0])
	//rows, err := db.Where(conditionString+condition+" ?", conditionValue).Find(&movies).Rows()
	var rows *sql.Rows
	if conditionString == "rating"{
		val, _ := strconv.Atoi(conditionValue)
		rows, err = db.Model(&TblImdb{}).Where(conditionString+condition+" ?",val ).Select("*").Rows()
	}else{
		rows, err = db.Model(&TblImdb{}).Where(conditionString+condition+" ?",conditionValue ).Select("*").Rows()
	}


	if err != nil {
		imdblog.WriteFile(FileObj, "failed to fetch values", os.Args[0])
		panic("failed to connect database")
	}
	defer rows.Close()
	var moviesList []*TblImdb
	for rows.Next() {
		// ScanRows scan a row into user
		err = db.ScanRows(rows, &movies)
		if err != nil {
			imdblog.WriteFile(FileObj, "failed to fetch values", os.Args[0])
			panic("failed to connect database")
		}
		moviesList = append(moviesList,&movies)
	}
	if moviesList != nil{
		return moviesList, true

	}
	return nil, false
}

func SearchInImdb(title string)  *TblImdb {

	key := imdb.Init("7283f93a")
	query := imdb.QueryData{Title: title}
	result , err := key.MovieByTitle(&query)

	if err != nil {
		imdblog.WriteFile(FileObj, "failed to connect database", os.Args[0])
		panic("failed to connect database")
	}
	rating, err := strconv.ParseFloat(result.ImdbRating,64)
	if err != nil {
		imdblog.WriteFile(FileObj, "failed to connect database", os.Args[0])
		panic("failed to connect database")
	}
	tblimdb := &TblImdb{Title: result.Title,
						ReleasedYear: result.Year,
						Rating: int(rating),
						Genres: result.Genre,

	}
	imdblog.WriteFile(FileObj, "TblImdb connected database", os.Args[0])
	vals := SqliteCreateData(tblimdb)
	return &vals
}