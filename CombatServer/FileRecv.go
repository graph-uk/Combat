package main

import (
	"crypto/md5"

	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func uploadSessionHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)

		sessionName := r.FormValue("SessionName")
		if sessionName == "" {
			fmt.Println("cannot extract session name")
			return
		}
		os.MkdirAll("./sessions/"+sessionName, 0777)
		f, err := os.OpenFile("./sessions/"+sessionName+"/"+filepath.Base(handler.Filename), os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		db, err := sql.Open("sqlite3", "./base.sl3")
		check(err)
		defer db.Close()

		req, err := db.Prepare("INSERT INTO Sessions(id,filename) VALUES(?,?)")
		check(err)
		_, err = req.Exec(sessionName, "astaxie")
		check(err)

		fmt.Println(r.Host + " Create new session: " + r.FormValue("SessionName"))
	}

}
