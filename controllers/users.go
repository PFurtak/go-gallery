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

// Create is used to create a new user account from signup form
// POST /signup
func (u *Users) Create(rw http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(rw, form)

}

// New is used to render the signup form for users to create an account.
// GET /signup
func (u *Users) New(rw http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(rw, nil); err != nil {
		panic(err)
	}
}

type Users struct {
	NewView *views.View
}

type SignUpForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
