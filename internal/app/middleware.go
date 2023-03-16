package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/with-insomnia/Hotel/pkg"
)

var AuthPaths = make(map[string]struct{})

type keyUserType string

const keyUser = "user"

func AddAuthPaths(paths ...string) {
	for _, path := range paths {
		AuthPaths[path] = struct{}{}
	}
}

func (app *App) authorizedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := AuthPaths[r.URL.Path]; !ok {
			pkg.ErrorHandler(w, http.StatusNotFound)
			return
		}
		cookie, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println(1)
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		session, err := app.sessionService.GetSessionByToken(cookie.Value)
		if err != nil {
			fmt.Println(err)
			fmt.Println(2)
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}
		if session.Expiry.Before(time.Now().UTC()) {
			fmt.Println(3)
			fmt.Println(session.Expiry)
			fmt.Println(time.Now())
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}

		user, err := app.userService.GetUserByToken(cookie.Value)
		if err != nil {
			fmt.Println(err)
			fmt.Println(4)
			http.Redirect(w, r, "/sign-in", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), keyUserType(keyUser), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) nonAuthorizedMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := AuthPaths[r.URL.Path]; !ok {
			pkg.ErrorHandler(w, http.StatusNotFound)
			return
		}
		c, err := r.Cookie("session_token")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		checkSessionFromDb, err := app.sessionService.GetSessionByToken(c.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		if checkSessionFromDb.Expiry.Before(time.Now()) {
			next.ServeHTTP(w, r)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	})
}
