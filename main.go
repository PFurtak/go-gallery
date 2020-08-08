package main

import (
	"fmt"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
)

var (
	homeTemplate    *template.Template
	contactTemplate *template.Template
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	if err := homeTemplate.Execute(rw, nil); err != nil {
		panic(err)
	}
}

func contactHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	if err := contactTemplate.Execute(rw, nil); err != nil {
		panic(err)
	}
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
	var err error

	homeTemplate, err = template.ParseFiles("views/home.gohtml", "views/layouts/footer.gohtml")
	if err != nil {
		panic(err)
	}

	contactTemplate, err = template.ParseFiles("views/contact.gohtml", "views/layouts/footer.gohtml")
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/", homeHandler)
	router.HandleFunc("/contact", contactHandler)
	router.HandleFunc("/faq", faqHandler)
	http.ListenAndServe(":5000", router)
}
