package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/pkg"
)

const (
	Layout = "2006-01-02"
)

var ErrNoRoom = errors.New("don't have free rooms")

var Alarm model.Notification

func (app *App) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		err := r.ParseForm()
		if err != nil {
			pkg.ErrorHandler(w, http.StatusInternalServerError)
			return
		}
		checkIn := r.Form.Get("check_in")
		checkOut := r.Form.Get("check_out")
		var room string
		rooms, err := strconv.Atoi(r.Form.Get("rooms"))
		if err != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		if rooms < 1 || rooms > 4 {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		switch rooms {
		case 1:
			room = "generalsolo"
		case 2:
			room = "generalduo"
		case 3:
			room = "luxduo"
		case 4:
			room = "luxtrio"
		}

		ok := pkg.CheckValidData(checkIn, checkOut)
		if !ok {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}

		freeRoom, err := app.roomService.CheckFreeRooms(room)
		if err != nil {
			Alarm.Message = fmt.Sprintf("room occupied")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		if freeRoom.Room_id == 0 {
			Alarm.Message = fmt.Sprintf("room occupied")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		Alarm.Message = fmt.Sprintf("check-in: [%s] check-out: [%s] rooms: [%s]", checkIn, checkOut, freeRoom.Room_name)
		Alarm.CheckIn = checkIn
		Alarm.CheckOut = checkOut
		Alarm.Room = room
		http.Redirect(w, r, "/make-reservation", http.StatusFound)
		return
	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}

func (app *App) Reservation(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		pkg.RenderTemplate(w, "make-reservation.html", Alarm)
		pkg.ClearStruct(&Alarm)
		return
	case http.MethodPost:
		r.ParseForm()
		firstname := r.FormValue("first-name")
		lastname := r.FormValue("last-name")
		email := r.FormValue("email")
		roomName := r.FormValue("room")
		checkin := r.FormValue("check-in")
		checkout := r.FormValue("check-out")

		fmt.Println(firstname, lastname, email, checkin, checkout, roomName)

		room, err := app.roomService.CheckFreeRooms(roomName)
		if err != nil {
			Alarm.Message = fmt.Sprintf("Don't have free %s room for date check-in: %s and check-out: %s", roomName, checkin, checkout)
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		user, ok := r.Context().Value(keyUserType(keyUser)).(model.User)
		if !ok {
			pkg.ErrorHandler(w, http.StatusUnauthorized)
			return
		}

		// if we have free roooms need logic to occupied room
		inTime, _ := time.Parse(Layout, checkin)
		outTime, _ := time.Parse(Layout, checkout)
		err = app.roomService.AddToWaitlist(model.Reservation{
			UserId:   user.ID,
			RoomID:   room.Room_id,
			CheckIn:  inTime,
			CheckOut: outTime,
		})
		if err != nil {
			pkg.ErrorHandler(w, http.StatusBadRequest)
			return
		}
		Alarm.Message = "Success"
		http.Redirect(w, r, "/", http.StatusFound)
		return
	default:
		pkg.ErrorHandler(w, http.StatusMethodNotAllowed)
		return
	}
}
