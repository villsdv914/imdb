package main

import (
	"Imdb/Sql"
	"Imdb/imdblog"
	"Imdb/utils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/gddo/httputil/header"
	"io"
	"log"
	"net/http"
	"os"
)

const logger = "imdb"

func checkError(err error) {
	if err != nil {
		stringError := err.Error()
		imdblog.WriteFile(logger, stringError, os.Args[0])
	}
}


func postHandler(w http.ResponseWriter, r *http.Request) {

	var jBody utils.RequestBody
	var searchBody utils.SearchBody
	var err error

	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return
		}
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()


	if r.RequestURI == "/api/search"{

		err = dec.Decode(&searchBody)
	}else{
		err = dec.Decode(&jBody)
	}

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {

		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			http.Error(w, msg, http.StatusBadRequest)


		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := fmt.Sprintf("Request body contains badly-formed JSON")
			http.Error(w, msg, http.StatusBadRequest)


		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			http.Error(w, msg, http.StatusBadRequest)




		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			http.Error(w, msg, http.StatusBadRequest)


		default:
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	tblMap := &sql.TblImdb{Title: jBody.Title,
		Rating:       jBody.Rating,
		ReleasedYear: jBody.ReleasedYear,
		Genres:       jBody.Genres,
	}

	if r.RequestURI == "/api/create" {

		tblData := sql.SqliteCreateData(tblMap)
		if tblData.ID > 0{
			jresp := utils.ResponseBody{Id: tblData.ID,
				Title: tblData.Title,
				Rating: tblData.Rating,
				ReleasedYear: tblData.ReleasedYear,
				Genres: tblData.Genres,
				CreatedAt: tblData.CreatedAt,
				UpdatedAt: tblData.UpdatedAt,


				}
			w.WriteHeader(http.StatusCreated)
			err := json.NewEncoder(w).Encode(jresp)
			checkError(err)
		}


	}else if r.RequestURI == "/api/edit" {
		conditon := &sql.TblImdb{Title :jBody.Title}
		if sql.SqliteUpdateData(sql.TblImdb{Genres: jBody.Genres, Rating: jBody.Rating}, conditon){
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(struct {
				Response string `json:"response,omitempty"`
			}{Response: "updated successfully"})
			checkError(err)
		}
	}else if r.RequestURI == "/api/search"{
		var vals interface{}
		vals, bvalue := sql.Search(searchBody.Condition,searchBody.SearchType, searchBody.TypeValue)
		if bvalue == false{
			vals = sql.SearchInImdb(searchBody.TypeValue)
		}

		fmt.Print(vals)
		err := json.NewEncoder(w).Encode(vals)
		checkError(err)
	}

}

func main() {
	imdblog.WriteFile(logger, "initializing", os.Args[0])
	sql.SqliteMigrate()
	mux := http.NewServeMux()
	mux.HandleFunc("/api/create", postHandler)
	mux.HandleFunc("/api/edit", postHandler)
	mux.HandleFunc("/api/search", postHandler)

	errServe := http.ListenAndServe(":1220", mux)
	checkError(errServe)
}
