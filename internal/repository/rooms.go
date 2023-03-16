package repository

import (
	"database/sql"
	"errors"

	"github.com/with-insomnia/Hotel/internal/model"
)

type RoomQuery interface {
	CheckFreeRooms(roomType string) (model.Rooms, error)
	Waitlist(reservation model.Reservation) error
	UpdateFreeRoom(roomId int) error
	GetAllWaitlist() ([]model.Reservation, error)
	GetRoomById(id int64) (model.Rooms, error)
	DeleteWaitList(roomId, userId int) error
}

type roomQuery struct {
	db *sql.DB
}

func (r *roomQuery) GetAllWaitlist() ([]model.Reservation, error) {
	rows, err := r.db.Query("SELECT * FROM waitlist")
	if err != nil {
		return []model.Reservation{}, err
	}
	defer rows.Close()
	var all []model.Reservation
	for rows.Next() {
		var post model.Reservation
		if err := rows.Scan(&post.UserId, &post.RoomID, &post.CheckIn, &post.CheckOut); err != nil {
			return []model.Reservation{}, err
		}
		all = append(all, post)
	}
	return all, nil
}

func (r *roomQuery) CheckFreeRooms(roomType string) (model.Rooms, error) {
	query := `SELECT * FROM rooms WHERE(room_name = $1 AND room_occupied = 0)`
	var result model.Rooms
	err := r.db.QueryRow(query, roomType).Scan(&result.Room_id, &result.Room_name, &result.Room_occupied)
	return result, err
}

func (r *roomQuery) GetRoomById(id int64) (model.Rooms, error) {
	query := `SELECT * FROM rooms WHERE(room_id = $1)`
	var result model.Rooms
	err := r.db.QueryRow(query, id).Scan(&result.Room_id, &result.Room_name, &result.Room_occupied)
	return result, err
}

func (r *roomQuery) Waitlist(reservation model.Reservation) error {
	query := `INSERT INTO waitlist(user_id, room_id, check_in, check_out) VALUES($1, $2, $3, $4)`
	_, err := r.db.Exec(query, reservation.UserId, reservation.RoomID, reservation.CheckIn, reservation.CheckOut)
	return err
}

func (r *roomQuery) UpdateFreeRoom(roomId int) error {
	query := `UPDATE rooms
	SET room_occupied = 1
	WHERE room_id = $1;`

	_, err := r.db.Exec(query, roomId)
	return err
}

func (r *roomQuery) DeleteWaitList(roomId, userId int) error {
	query := `DELETE FROM waitlist WHERE room_id = $1 AND user_id = $2`
	res, err := r.db.Exec(query, roomId, userId)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("error delete from waitlist")
	}
	return nil
}
