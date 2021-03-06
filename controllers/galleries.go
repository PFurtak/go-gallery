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
	ShowGallery     = "show_gallery"
	EditGallery     = "edit_gallery"
	maxMultipartmem = 1 << 20 //1mb
)

func NewGalleries(gs models.GalleryService, is models.ImageService, router *mux.Router) *Galleries {
	return &Galleries{
		New:       views.NewView("bootstrap", "galleries/new"),
		ShowView:  views.NewView("bootstrap", "galleries/show"),
		EditView:  views.NewView("bootstrap", "galleries/edit"),
		IndexView: views.NewView("bootstrap", "galleries/index"),
		gs:        gs,
		is:        is,
		router:    router,
	}
}

type Galleries struct {
	New       *views.View
	ShowView  *views.View
	EditView  *views.View
	IndexView *views.View
	gs        models.GalleryService
	is        models.ImageService
	router    *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

// GET /galleries/
func (g *Galleries) Index(rw http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.gs.ByUserID(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Something went wrong", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = galleries
	g.IndexView.Render(rw, r, vd)
}

// GET /galleries/:id
func (g *Galleries) Show(rw http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(rw, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Yield = gallery
	g.ShowView.Render(rw, r, vd)
}

// GET /galleries/:id/edit
func (g *Galleries) Edit(rw http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(rw, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(rw, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	vd.User = user
	g.EditView.Render(rw, r, vd)
}

// POST /galleries/:id/update
func (g *Galleries) Update(rw http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(rw, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(rw, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.EditView.Render(rw, r, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(rw, r, vd)
		return
	}
	vd.Alert = &views.Alert{
		Level:     views.AlertLvlSuccess,
		AlertType: views.AlertTypeSuccess,
		Message:   "Gallery successfully updated!",
	}
	g.EditView.Render(rw, r, vd)
}

// Create POSTS /galleries
func (g *Galleries) Create(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		g.New.Render(rw, r, vd)
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
		g.New.Render(rw, r, vd)
		return
	}
	url, err := g.router.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(rw, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(rw, r, url.Path, http.StatusFound)
}

func (g *Galleries) galleryByID(rw http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(rw, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(rw, "Gallery not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(rw, "Something went wrong ;[", http.StatusInternalServerError)
		}
		return nil, err
	}
	images, _ := g.is.ByGalleryID(gallery.ID)
	gallery.Images = images
	return gallery, nil
}

// POST /galleries/:id/delete
func (g *Galleries) Delete(rw http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(rw, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(rw, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = gallery
		g.EditView.Render(rw, r, vd)
		return
	}
	http.Redirect(rw, r, "/galleries", http.StatusFound)
}

// POST /galleries/:id/images
func (g *Galleries) ImageUpload(rw http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(rw, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(rw, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Yield = gallery
	err = r.ParseMultipartForm(maxMultipartmem)
	if err != nil {
		vd.SetAlert(err)
		g.EditView.Render(rw, r, vd)
		return
	}

	files := r.MultipartForm.File["images"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(rw, r, vd)
			return
		}
		defer file.Close()
		err = g.is.Create(gallery.ID, file, f.Filename)
		if err != nil {
			vd.SetAlert(err)
			g.EditView.Render(rw, r, vd)
			return
		}
	}
	url, err := g.router.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(rw, r, "/galleries", http.StatusFound)
	}
	http.Redirect(rw, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/images/:filename/delete
func (g *Galleries) ImageDelete(rw http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(rw, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(rw, "Gallery not found", http.StatusNotFound)
		return
	}
	filename := mux.Vars(r)["filename"]
	i := models.Image{
		Filename:  filename,
		GalleryID: gallery.ID,
	}
	err = g.is.Delete(&i)
	if err != nil {
		var vd views.Data
		vd.Yield = gallery
		vd.SetAlert(err)
		g.EditView.Render(rw, r, vd)
		return
	}
	url, err := g.router.Get(EditGallery).URL("id", fmt.Sprintf("%v", gallery.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(rw, r, "/galleries", http.StatusFound)
		return
	}
	http.Redirect(rw, r, url.Path, http.StatusFound)
}
