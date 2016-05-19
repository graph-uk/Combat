package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	//	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func markCaseFailed(caseID string) {
	db, err := sql.Open("sqlite3", "./base.sl3")
	check(err)
	defer db.Close()
	req, err := db.Prepare(`UPDATE Cases SET inProgress="false", passed="false", finished="true" WHERE id=?`)
	check(err)
	_, err = req.Exec(caseID)
	check(err)
}

func markCasePassed(caseID string) {
	db, err := sql.Open("sqlite3", "./base.sl3")
	check(err)
	defer db.Close()
	req, err := db.Prepare(`UPDATE Cases SET inProgress="false", passed="true", finished="true" WHERE id=?`)
	check(err)
	_, err = req.Exec(caseID)
	check(err)
}

func markCaseNotInProgress(caseID string) {
	db, err := sql.Open("sqlite3", "./base.sl3")
	check(err)
	defer db.Close()
	req, err := db.Prepare(`UPDATE Cases SET inProgress="false" WHERE id=?`)
	check(err)
	_, err = req.Exec(caseID)
	check(err)
}

func setCaseResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

	} else {
		//tryName := strconv.FormatInt(time.Now().UnixNano(), 10)
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		caseID := r.FormValue("caseID")
		if caseID == "" {
			fmt.Println("cannot extract caseID")
			return
		}

		exitStatus := r.FormValue("exitStatus")
		if exitStatus == "" {
			fmt.Println("cannot extract exitStatus")
			return
		}

		stdOut := r.FormValue("stdOut")
		if stdOut == "" {
			fmt.Println("cannot extract stdOut")
			return
		}

		db, err := sql.Open("sqlite3", "./base.sl3")
		check(err)
		defer db.Close()

		req, err := db.Prepare(`SELECT id FROM Tries WHERE caseID=?`)
		check(err)
		rows, err := req.Query(caseID)
		check(err)

		triesCount := 0
		for rows.Next() {
			triesCount++
		}
		rows.Close()

		fmt.Println("CurrentTryCount=" + strconv.Itoa(triesCount))

		req, err = db.Prepare("INSERT INTO Tries(caseID,exitStatus,stdOut) VALUES(?,?,?)")
		check(err)
		res, err := req.Exec(caseID, exitStatus, stdOut)
		check(err)
		tryID64, err := res.LastInsertId()
		check(err)
		tryID := strconv.Itoa(int(tryID64))
		//fmt.Println(strconv.Itoa(int(tryID)))

		db.Close()
		if triesCount > 2 && exitStatus != "0" {
			markCaseFailed(caseID)
		} else {
			if exitStatus == "0" {
				markCasePassed(caseID)
			} else {
				markCaseNotInProgress(caseID)
			}
		}

		os.MkdirAll("./tries/"+tryID, 0777)
		f, err := os.OpenFile("./tries/"+tryID+"/out_archived.zip", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		fmt.Println(r.Host + " ran case: " + caseID)
	}
}
