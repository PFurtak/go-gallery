package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Users/patrickfurtak/desktop/go-gallery/context"
	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
	"github.com/Users/patrickfurtak/desktop/go-gallery/views"
	"github.com/gorilla/mux"
)

const (
	ShowGallery = "show_gallery"
)

func NewGalleries(gs models.GalleryService, router *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		gs:       gs,
		router:   router,
	}
}

type Galleries struct {
	New      *views.View
	ShowView *views.View
	gs       models.GalleryService
	router   *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries/:id
func (g *Galleries) Show(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(rw, "Invalid gallery ID", http.StatusNotFound)
		return
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(rw, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(rw, "Something went wrong ;[", http.StatusInternalServerError)
		}
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(rw, vd)
}

// Create POSTS /galleries
func (g *Galleries) Create(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(rw, vd)
		return
	}
	user := context.User(r.Context())
	if user == nil {
		http.Redirect(rw, r, "/login", http.StatusFound)
	}
	fmt.Println("Create user: ", user)
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}

	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(rw, vd)
		return
	}
	url, err := g.router.Get(ShowGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		// TODO make this go to index
		http.Redirect(rw, r, "/", http.StatusFound)
		return
	}
	http.Redirect(rw, r, url.Path, http.StatusFound)
	fmt.Fprintln(rw, gallery)
}
