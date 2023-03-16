package app

import (
	"net/http"
	"time"

	"github.com/with-insomnia/Hotel/internal/config"
)

func (app *App) Routes(cfg config.Http) *http.Server {
	authPaths := []string{
		"/",
		"/post",
		"/welcome",
		"/sign-in",
		"/sign-up",
		"/logout",
		"/sweetalert2.min.css",
		"/sweetalert2.all.min.js",
		"/check-availability",
		"/make-reservation",
		"/sweetalert2.all.min.js/",
		"/admin",
		"/query-list",
		"/dashboard",
		"/rooms",
		"/message-list",
	}

	AddAuthPaths(authPaths...)

	mux := http.NewServeMux()

	mux.HandleFunc("/", app.authorizedMiddleware(app.IndexHandler))
	mux.HandleFunc("/post", app.authorizedMiddleware(app.PostHandler))
	mux.HandleFunc("/welcome", app.nonAuthorizedMiddleware(app.WelcomeHandler))

	mux.HandleFunc("/check-availability", app.authorizedMiddleware(app.CheckAvailability))
	mux.HandleFunc("/make-reservation", app.authorizedMiddleware(app.Reservation))

	mux.HandleFunc("/sign-in", app.nonAuthorizedMiddleware(app.LoginHandler))
	mux.HandleFunc("/sign-up", app.nonAuthorizedMiddleware(app.RegisterHandler))
	mux.HandleFunc("/logout", app.authorizedMiddleware(app.LogoutHandler))

	mux.HandleFunc("/dashboard", app.authorizedMiddleware(app.Admin))
	mux.HandleFunc("/query-list", app.authorizedMiddleware(app.AdminQuery))
	mux.HandleFunc("/message-list", app.authorizedMiddleware(app.AdminMessages))
	mux.HandleFunc("/rooms", app.authorizedMiddleware(app.AdminRooms))

	fs := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	server := &http.Server{
		ReadTimeout:  time.Second * time.Duration(cfg.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.IdleTimeout),
		Addr:         cfg.Port,
		Handler:      mux,
	}
	return server
}
