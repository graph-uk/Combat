package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	fmt.Printf(r.UserAgent())
	fmt.Println("*")
	//fmt.Printf(r.)
	fmt.Println("*")
}

func main() {
	http.HandleFunc("/uploadSession", handler)
	http.ListenAndServe(":9090", nil)
}
