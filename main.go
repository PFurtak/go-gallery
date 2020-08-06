package main

import (
	"fmt"
	"net/http"
)

func handlerFunc(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	if r.URL.Path == "/" {
		fmt.Fprint(rw, "<h1>Welcome to Go-Gallery!</h1>")
	} else {
		fmt.Fprint(rw, "To get in touch, please email: <a href=\"mailto:support@gogallery.com\">support@gogallery.com</a>")
	}
}

func main() {
	http.HandleFunc("/", handlerFunc)
	http.ListenAndServe(":5000", nil)

}
