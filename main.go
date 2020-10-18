package main

import (
	"fmt"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/controllers"
	"github.com/Users/patrickfurtak/desktop/go-gallery/middleware"
	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
	"github.com/Users/patrickfurtak/desktop/go-gallery/rand"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
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
	cfg := DefaultConfig()
	dbCfg := DefaultPostgresConfig()
	services, err := models.NewServices(dbCfg.Dialect(), dbCfg.DbConfigInfo())
	must(err)
	defer services.Close()
	services.AutoMigrate()
	// services.DestructiveReset()

	router := mux.NewRouter()
	staticController := controllers.NewStatic()
	usersController := controllers.NewUsers(services.User)
	galleriesController := controllers.NewGalleries(services.Gallery, services.Image, router)

	b, err := rand.Bytes(32)
	must(err)

	csrfMw := csrf.Protect(b, csrf.Secure(cfg.IsProd()))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{User: userMw}

	router.NotFoundHandler = http.HandlerFunc(notFound)
	router.Handle("/", staticController.Home).Methods("GET")
	router.Handle("/contact", staticController.Contact).Methods("GET")
	router.HandleFunc("/signup", usersController.New).Methods("GET")
	router.HandleFunc("/signup", usersController.Create).Methods("POST")
	router.Handle("/login", usersController.LoginView).Methods("GET")
	router.HandleFunc("/login", usersController.Login).Methods("POST")
	router.HandleFunc("/faq", faqHandler)

	//Assets
	assetHandler := http.FileServer(http.Dir("./assets"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	router.PathPrefix("/assets").Handler(assetHandler)

	//Images routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	router.PathPrefix("/images/").Handler(http.StripPrefix("/images", imageHandler))

	//Gallery routes

	router.Handle("/galleries", requireUserMw.Applyfn(galleriesController.Index)).Methods("GET")
	router.Handle("/galleries/new", requireUserMw.Apply(galleriesController.New)).Methods("GET")
	router.HandleFunc("/galleries", requireUserMw.Applyfn(galleriesController.Create)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}", galleriesController.Show).Methods("GET").Name(controllers.ShowGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.Applyfn(galleriesController.Edit)).Methods("GET").Name(controllers.EditGallery)
	router.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.Applyfn(galleriesController.Update)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.Applyfn(galleriesController.Delete)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.Applyfn(galleriesController.ImageUpload)).Methods("POST")
	router.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.Applyfn(galleriesController.ImageDelete)).Methods("POST")

	http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), csrfMw(userMw.Apply(router)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
