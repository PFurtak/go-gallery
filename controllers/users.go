package controllers

import (
	"fmt"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
	"github.com/Users/patrickfurtak/desktop/go-gallery/views"
)

// NewUsers is used to create a new Users controller
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/newusers"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// Create is used to create a new user account from signup form
// POST /signup
func (u *Users) Create(rw http.ResponseWriter, r *http.Request) {
	var form SignUpForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(rw, user)
}

// New is used to render the signup form for users to create an account.
// GET /signup
func (u *Users) New(rw http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(rw, nil); err != nil {
		panic(err)
	}
}

// Login is used to verify provided credentials and login if correct
// POST /login
func (u *Users) Login(rw http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(rw, form)
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

type SignUpForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
