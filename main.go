package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {

	type User struct {
		Name string
	}

	data := User{
		Name: "John Smith",
	}

	rw.Header().Set("Content-Type", "text/html")
	t, err := template.ParseFiles("templates/hello.gohtml")
	if err != nil {
		panic(err)
	}

	err = t.Execute(os.Stdout, data)

	if err != nil {
		panic(err)
	}

	data.Name = "Patrick F."
	err = t.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}

func contactHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, "To get in touch, please email: <a href=\"mailto:support@gogallery.com\">support@gogallery.com</a>")
}

func faqHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	fmt.Fprint(rw, "<h1>FAQ</h1><br><ol><li>Who is this site for? <b>Photographers!</b></li><li>Can I upload my own photos? <b>Yes!</b></li><li>What language is this application written in? <b>Golang</b></li></ol>")
}

func notFound(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-type", "text/html")
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprint(rw, "<h1>404</h1><br><h3>Page not found :[</h3>")
}

func main() {
	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/contact", contactHandler)
	router.HandleFunc("/faq", faqHandler)
	http.ListenAndServe(":5000", router)
}
