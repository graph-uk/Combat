package main

import (
	"fmt"
	//"io"
	"io/ioutil"
	"net/http"
	//"time"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func sendJobToNode(w http.ResponseWriter, r *http.Request) {
}

func getJobHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		db, err := sql.Open("sqlite3", "./base.sl3")
		check(err)
		defer db.Close()
		rows, err := db.Query("SELECT id,params FROM Sessions WHERE status=0 limit 1")
		check(err)

		var sessionId, sessionParams string
		if rows.Next() {
			err = rows.Scan(&sessionId, &sessionParams)
			check(err)
			rows.Close()

			req, err := db.Prepare("UPDATE Sessions SET status=? WHERE id=?")
			check(err)
			_, err = req.Exec(1, sessionId)
			check(err)

			w.Header().Add("Command", "CasesExplore")
			w.Header().Add("Params", sessionParams)
			zipArchive, err := ioutil.ReadFile("./sessions/" + sessionId + "/archived.zip")
			check(err)
			_, err = w.Write(zipArchive)
			check(err)
			fmt.Println(r.Host + " Get a job: " + sessionId)
		} else {
			w.Header().Add("Command", "idle")
		}
	}
}
