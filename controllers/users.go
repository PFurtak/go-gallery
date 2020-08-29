package controllers

import (
	"log"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
	"github.com/Users/patrickfurtak/desktop/go-gallery/rand"
	"github.com/Users/patrickfurtak/desktop/go-gallery/views"
)

// NewUsers is used to create a new Users controller
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/newusers"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// Create is used to create a new user account from signup form
// POST /signup
func (u *Users) Create(rw http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form SignUpForm
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.NewView.Render(rw, r, vd)
		return
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}

	if err := u.us.Create(&user); err != nil {
		vd.SetAlert(err)
		u.NewView.Render(rw, r, vd)
		return
	}
	err := u.signIn(rw, &user)
	if err != nil {
		http.Redirect(rw, r, "/login", http.StatusFound)
		return
	}
	http.Redirect(rw, r, "/galleries", http.StatusFound)
}

// New is used to render the signup form for users to create an account.
// GET /signup
func (u *Users) New(rw http.ResponseWriter, r *http.Request) {
	u.NewView.Render(rw, r, nil)
}

// Login is used to parse login form on submit
// POST /login
func (u *Users) Login(rw http.ResponseWriter, r *http.Request) {
	vd := views.Data{}
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		log.Println(err)
		vd.SetAlert(err)
		u.LoginView.Render(rw, r, vd)
		return
	}

	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			vd.Alert = &views.Alert{
				Level:     views.AlertLvlError,
				AlertType: views.AlertTypeError,
				Message:   "Invalid email address",
			}
		default:
			vd.SetAlert(err)
		}
		u.LoginView.Render(rw, r, vd)
		return
	}
	err = u.signIn(rw, user)
	if err != nil {
		vd.SetAlert(err)
		u.LoginView.Render(rw, r, vd)
		return
	}
	http.Redirect(rw, r, "/galleries", http.StatusFound)
}

func (u *Users) signIn(rw http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}

	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
	}
	http.SetCookie(rw, &cookie)
	return nil
}

type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
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
