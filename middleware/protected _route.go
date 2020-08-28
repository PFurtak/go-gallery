package middleware

import (
	"fmt"
	"net/http"

	"github.com/Users/patrickfurtak/desktop/go-gallery/context"
	"github.com/Users/patrickfurtak/desktop/go-gallery/models"
)

type User struct {
	models.UserService
}

func (mw *User) Apply(next http.Handler) http.HandlerFunc {
	return mw.Applyfn(next.ServeHTTP)
}

func (mw *User) Applyfn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("remember_token")
		if err != nil {
			next(rw, r)
			return
		}
		user, err := mw.UserService.ByRemember(cookie.Value)
		if err != nil {
			next(rw, r)
			return
		}
		ctx := r.Context()
		ctx = context.WithUser(ctx, user)
		r = r.WithContext(ctx)
		fmt.Println("User found: ", user)
		next(rw, r)
	})
}

type RequireUser struct {
	User
}

func (mw *RequireUser) Apply(next http.Handler) http.HandlerFunc {
	return mw.Applyfn(next.ServeHTTP)
}

func (mw *RequireUser) Applyfn(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := context.User(r.Context())
		if user == nil {
			http.Redirect(rw, r, "/login", http.StatusFound)
			return
		}
		next(rw, r)
	})
}
