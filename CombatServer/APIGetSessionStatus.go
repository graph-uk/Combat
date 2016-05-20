package main

import (
	"fmt"
	//"io"
	//"io/ioutil"
	"net/http"
	//"time"
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func getSessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//sessionID := r.Header.Get("sessionID")
		r.ParseMultipartForm(32 << 20)
		sessionID := r.FormValue("sessionID")
		if sessionID == "" {
			fmt.Println("cannot extract session ID")
			return
		}

		db, err := sql.Open("sqlite3", "./base.sl3")
		check(err)
		defer db.Close()
		req, err := db.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=?`)
		check(err)
		rows, err := req.Query(sessionID)
		check(err)
		var totalCasesCount int
		rows.Next()
		err = rows.Scan(&totalCasesCount)
		check(err)
		rows.Close()

		req, err = db.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=? AND finished="true"`)
		check(err)
		rows, err = req.Query(sessionID)
		check(err)
		var finishedCasesCount int
		rows.Next()
		err = rows.Scan(&finishedCasesCount)
		check(err)
		rows.Close()

		req, err = db.Prepare(`SELECT Count()as count FROM Cases WHERE sessionID=? AND finished="true" AND passed="false"`)
		check(err)
		rows, err = req.Query(sessionID)
		check(err)
		var failedCasesCount int
		rows.Next()
		err = rows.Scan(&failedCasesCount)
		check(err)
		rows.Close()

		req, err = db.Prepare(`SELECT cmdLine FROM Cases WHERE sessionID=? AND finished="true" AND passed="false"`)
		check(err)
		rows, err = req.Query(sessionID)
		check(err)
		var errorCases []string
		for rows.Next() {
			var cmdLine string
			err = rows.Scan(&cmdLine)
			check(err)
			errorCases = append(errorCases, cmdLine)
		}
		rows.Close()
		//rows.Next()
		//err = rows.Scan(&failedCasesCount)
		//check(err)

		if totalCasesCount == finishedCasesCount && totalCasesCount != 0 {
			w.Header().Set("Finished", "True")
			if failedCasesCount == 0 {
				w.Write([]byte("Success. All tests passed"))
			} else {
				w.Write([]byte("Finished. Errors: " + strconv.Itoa(failedCasesCount) + "\r\n"))
				for _, curCase := range errorCases {
					w.Write([]byte("    " + curCase + "\r\n"))
				}
			}
		} else {
			w.Header().Set("Finished", "False")
			if totalCasesCount != 0 {
				if failedCasesCount != 0 {
					w.Write([]byte("Running (" + strconv.Itoa(finishedCasesCount) + "/" + strconv.Itoa(totalCasesCount) + ") Errors: " + strconv.Itoa(failedCasesCount) + "\r\n"))
					for _, curCase := range errorCases {
						w.Write([]byte("    " + curCase + "\r\n"))
					}
				} else {
					w.Write([]byte("Running (" + strconv.Itoa(finishedCasesCount) + "/" + strconv.Itoa(totalCasesCount) + ")"))
				}
			} else {
				w.Write([]byte("Cases exploring"))
			}
		}

		fmt.Println(r.Host + " Get session status: for session: " + sessionID)
	}
}
