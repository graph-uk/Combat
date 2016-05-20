package main

import (
	"fmt"
	//"io"
	"net/http"
	//"os"
	"strconv"
	"strings"
	//"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func setSessionCasesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		r.ParseMultipartForm(32 << 20)

		sessionID := r.FormValue("sessionID")
		if sessionID == "" {
			fmt.Println("cannot extract session ID")
			return
		}

		sessionCases := r.FormValue("cases")
		if sessionCases == "" {
			fmt.Println("cannot extract session cases")
			return
		}

		sessionCasesArr := strings.Split(sessionCases, "\n")

		db, err := sql.Open("sqlite3", "./base.sl3")
		check(err)
		defer db.Close()

		req, err := db.Prepare("INSERT INTO Cases(cmdline, sessionID) VALUES(?,?)")
		check(err)

		casesCount := 0
		for _, curCase := range sessionCasesArr {
			curCaseCleared := strings.TrimSpace(curCase)
			if curCaseCleared != "" {
				casesCount++
				_, err = req.Exec(curCase, sessionID)
				check(err)
			}
		}

		//fmt.Println(r.Host + " Create new session: " + sessionName + " " + sessionParams)
		fmt.Println(r.Host + " Provided " + strconv.Itoa(casesCount) + " cases for session: " + sessionID)

	}
}
