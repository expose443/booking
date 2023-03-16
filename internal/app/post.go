package app

import (
	"log"
	"net/http"

	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/pkg"
)

func (app *App) PostHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		r.ParseForm()
		message := r.Form.Get("msg")
		if len(message) == 0 {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
		if !ok {
			pkg.ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		post := model.Post{
			Content: message,
			Author:  user,
		}

		status, err := app.postService.CreateAdminPost(&post)
		if err != nil {
			log.Println(err)
			switch status {
			case http.StatusInternalServerError:
				pkg.ErrorHandler(w, http.StatusInternalServerError)
				return
			case http.StatusBadRequest:
				pkg.ErrorHandler(w, http.StatusBadRequest)
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)

	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
