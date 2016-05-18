package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		sessionName := strconv.FormatInt(time.Now().UnixNano(), 10)
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		sessionParams := r.FormValue("params")
		if sessionParams == "" {
			fmt.Println("cannot extract session params")
			return
		}

		os.MkdirAll("./sessions/"+sessionName, 0777)
		f, err := os.OpenFile("./sessions/"+sessionName+"/archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		db, err := sql.Open("sqlite3", "./base.sl3")
		check(err)
		defer db.Close()

		req, err := db.Prepare("INSERT INTO Sessions(id,params) VALUES(?,?)")
		check(err)
		_, err = req.Exec(sessionName, sessionParams)
		check(err)

		io.WriteString(w, sessionName)
		fmt.Println(r.Host + " Create new session: " + sessionName + " " + sessionParams)
	}

}
