package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, "<h1>Welcome to Go-Gallery!</h1>")
}

func contactHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, "To get in touch, please email: <a href=\"mailto:support@gogallery.com\">support@gogallery.com</a>")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/contact", contactHandler)
	http.ListenAndServe(":5000", router)
}
