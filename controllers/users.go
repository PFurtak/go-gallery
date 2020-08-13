package controllers

import (
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/views"
)

// NewUsers is used to create a new Users controller

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/newusers.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

func (u *Users) New(rw http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(rw, nil); err != nil {
		panic(err)
	}
}
