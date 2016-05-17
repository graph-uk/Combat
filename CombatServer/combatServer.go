package main

import (
	"net/http"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

//var db *sql.DB

func main() {

	http.HandleFunc("/uploadSession", uploadSessionHandler)
	http.ListenAndServe(":9090", nil)
}
