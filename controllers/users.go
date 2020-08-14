package controllers

import (
	"fmt"
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

// New is used to render the signup form for users to create an account.
// GET /signup
func (u *Users) New(rw http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(rw, nil); err != nil {
		panic(err)
	}
}

// Create is used to create a new user account from signup form
// POST /signup
func (u *Users) Create(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "temp res")
}
