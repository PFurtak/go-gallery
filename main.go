package main

import (
	"fmt"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/controllers"
	"github.com/Users/patrickfurtak/desktop/go-gallery/views"
	"github.com/gorilla/mux"
)

var (
	homeView    *views.View
	contactView *views.View
)

func homeHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	must(homeView.Render(rw, nil))
}

func contactHandler(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	must(contactView.Render(rw, nil))
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
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	usersController := controllers.NewUsers()

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/contact", contactHandler).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	router.HandleFunc("/faq", faqHandler)
	http.ListenAndServe(":5000", router)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
