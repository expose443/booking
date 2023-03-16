package app

import (
	"log"
	"net/http"

	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/pkg"
)

func (app *App) IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		pkg.ErrorHandler(w, http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
	if !ok {
		pkg.ErrorHandler(w, http.StatusUnauthorized)
		return
	}

	post, err := app.postService.GetAllPosts()
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	for i, v := range post {
		author, err := app.userService.GetUserById(v.Author.ID)
		if err != nil {
			log.Println(err)
			pkg.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		post[i].Author.FirstName = author.FirstName
		post[i].Author.LastName = author.LastName

	}

	data := model.Data{
		Posts:        post,
		User:         user,
		Notification: Alarm,
	}
	pkg.RenderTemplate(w, "index.page.html", data)
	pkg.ClearStruct(&data)
	pkg.ClearStruct(&Alarm)
}

func (app *App) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
	post, err := app.postService.GetAllPosts()
	if err != nil {
		log.Println(err)
		pkg.ErrorHandler(w, http.StatusInternalServerError)
		return
	}
	data := model.Data{
		Posts:        post,
		Notification: Alarm,
	}
	pkg.RenderTemplate(w, "welcome.html", data)
	pkg.ClearStruct(&data.Notification)
	pkg.ClearStruct(&Alarm)
}
