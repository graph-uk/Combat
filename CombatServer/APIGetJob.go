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

		if rows.Next() {
			var sessionId, sessionParams string
			err = rows.Scan(&sessionId, &sessionParams)
			check(err)
			rows.Close()

			req, err := db.Prepare("UPDATE Sessions SET status=? WHERE id=?")
			check(err)
			_, err = req.Exec(1, sessionId)
			check(err)

			w.Header().Add("Command", "CasesExplore")
			w.Header().Add("Params", sessionParams)
			w.Header().Add("SessionID", sessionId)

			zipArchive, err := ioutil.ReadFile("./sessions/" + sessionId + "/archived.zip")
			check(err)
			_, err = w.Write(zipArchive)
			check(err)
			fmt.Println(r.Host + " Get a job (CasesExplore) for session: " + sessionId)
		} else { // when no one session needed to explore cases, check are we have cases to run
			var caseID, caseCMD, sessionID string
			rows, err := db.Query(`SELECT id, cmdLine, sessionID FROM cases WHERE finished="false" AND inProgress="false" ORDER BY RANDOM() LIMIT 1`)
			check(err)
			if rows.Next() {
				err = rows.Scan(&caseID, &caseCMD, &sessionID)
				check(err)
				rows.Close()

				req, err := db.Prepare("UPDATE Cases SET inProgress=? WHERE id=?")
				check(err)
				_, err = req.Exec(true, caseID)
				check(err)

				w.Header().Add("Command", "RunCase")
				w.Header().Add("Params", caseCMD)
				w.Header().Add("SessionID", caseID)

				zipArchive, err := ioutil.ReadFile("./sessions/" + sessionID + "/archived.zip")
				check(err)
				_, err = w.Write(zipArchive)
				check(err)
				fmt.Println(r.Host + " Get a job (CasesRun) for case: " + caseCMD)

			} else { // when not found cases to run
				w.Header().Add("Command", "idle")
			}
		}
	}
}
