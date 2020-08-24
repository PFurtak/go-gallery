package main

import (
	"fmt"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/controllers"
	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
	"github.com/gorilla/mux"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "patrickfurtak"
	dbname = "gogallery"
)

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

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable", host, port, user, dbname)
	services, err := models.NewServices(psqlInfo)
	must(err)
	// defer us.Close()
	// us.AutoMigrate()
	// us.DestructiveReset()

	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)

	router := mux.NewRouter()
	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")
	router.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")
	router.HandleFunc("/faq", faqHandler)
	http.ListenAndServe(":5000", router)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
