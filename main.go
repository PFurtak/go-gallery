package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, "<h1>Welcome to my cool site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":5000", nil)

}
