package app

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/pkg"
)

func (app *App) Admin(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
	if !ok {
		pkg.ErrorHandler(w, http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case http.MethodGet:

		pkg.RenderTemplate(w, "admin.html", nil)
		return
	case http.MethodPost:
		//
	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *App) AdminQuery(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
	if !ok {
		pkg.ErrorHandler(w, http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		waitlist, err := app.roomService.GetAllWaitlist()
		if err != nil {
			pkg.ErrorHandler(w, http.StatusInternalServerError)
			return
		}

		for i, v := range waitlist {
			author, err := app.userService.GetUserById(v.UserId)
			if err != nil {
				pkg.ErrorHandler(w, http.StatusInternalServerError)
				return
			}

			room, err := app.roomService.GetRoomById(v.RoomID)
			waitlist[i].Author = author
			waitlist[i].Room = room
			in := v.CheckIn.Format(Layout)
			out := v.CheckOut.Format(Layout)
			waitlist[i].In = in
			waitlist[i].Out = out

		}
		pkg.RenderTemplate(w, "admin-query.html", waitlist)
		return
	case http.MethodPost:
		r.ParseForm()
		room_id, err := strconv.Atoi(r.FormValue("room_id"))
		user_id, err1 := strconv.Atoi(r.FormValue("user_id"))
		if err != nil || err1 != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		if room_id < 1 || user_id < 1 {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		err = app.roomService.DeleteWaitList(room_id, user_id)
		if err != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/query-list", http.StatusFound)
		return

	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *App) AdminMessages(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
	if !ok {
		pkg.ErrorHandler(w, http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		post, err := app.postService.GetAllAdminPost()
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
			Notification: Alarm,
		}
		pkg.RenderTemplate(w, "admin-messages.html", data)
		pkg.ClearStruct(&data)
		return
	case http.MethodPost:
		r.ParseForm()
		id, err := strconv.Atoi(r.FormValue("id"))
		if err != nil || id < 1 {
			fmt.Println(1)
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		do := r.FormValue("do")
		fmt.Println(id)
		fmt.Println(do)
		switch do {
		case "public":
			post, err := app.postService.GetAdminPostById(int64(id))
			if err != nil {
				fmt.Println(2)
				pkg.ErrorHandler(w, http.StatusBadRequest)
				return
			}
			status, err := app.postService.CreatePost(&post)
			if status != 200 {
				pkg.ErrorHandler(w, status)
			}
			Alarm.Message = "Success"
			http.Redirect(w, r, "/message-list", http.StatusFound)
			return
		case "remove":
			err = app.postService.DeleteAdminPost(int64(id))
			if err != nil {
				fmt.Println(3)
				pkg.ErrorHandler(w, http.StatusBadRequest)
				return
			}
			Alarm.Message = "Success"
			http.Redirect(w, r, "/message-list", http.StatusFound)
			return
		default:
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return

		}

	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *App) AdminRooms(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
	if !ok {
		pkg.ErrorHandler(w, http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		pkg.RenderTemplate(w, "admin-rooms.html", nil)
		return
	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
