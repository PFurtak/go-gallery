package main

import (
	"fmt"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/controllers"
	"github.com/Users/patrickfurtak/desktop/go-gallery/middleware"
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
	defer services.Close()
	services.AutoMigrate()
	// services.DestructiveReset()

	router := mux.NewRouter()
	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery, router)

	requireUserMw := middleware.RequireUser{UserService: services.User}

	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")
	router.HandleFunc("/cookietest", usersController.CookieTest).Methods("GET")
	router.HandleFunc("/faq", faqHandler)

	//Gallery routes

	router.Handle("/galleries/new", requireUserMw.Apply(galleriesController.New)).Methods("GET")
	router.HandleFunc("/galleries", requireUserMw.Applyfn(galleriesController.Create)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleriesController.Show).Methods("GET").Name(controllers.ShowGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.Applyfn(galleriesController.Edit)).Methods("GET")
	router.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.Applyfn(galleriesController.Update)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.Applyfn(galleriesController.Delete)).Methods("POST")

	http.ListenAndServe(":5000", router)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
