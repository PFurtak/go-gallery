package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "<h1>Welcome to my site!</h1>")
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":5000", nil)

}
