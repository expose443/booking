package service

import (
	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/internal/repository"
)

type RoomService interface {
	CheckFreeRooms(roomType string) (model.Rooms, error)
	AddToWaitlist(reservation model.Reservation) error
	GetAllWaitlist() ([]model.Reservation, error)
	GetRoomById(id int64) (model.Rooms, error)
	DeleteWaitList(roomId, userId int) error
}

type roomService struct {
	repository.RoomQuery
}

func NewRoomService(dao repository.DAO) RoomService {
	return &roomService{
		dao.NewRoomQuery(),
	}
}

func (r *roomService) CheckFreeRooms(roomType string) (model.Rooms, error) {
	return r.RoomQuery.CheckFreeRooms(roomType)
}

func (r *roomService) AddToWaitlist(reservation model.Reservation) error {
	err := r.RoomQuery.Waitlist(reservation)
	if err != nil {
		return err
	}
	err = r.RoomQuery.UpdateFreeRoom(int(reservation.RoomID))
	return err
}

func (r *roomService) GetAllWaitlist() ([]model.Reservation, error) {
	waitlist, err := r.RoomQuery.GetAllWaitlist()
	if err != nil {
		return waitlist, err
	}
	return waitlist, nil
}

func (r *roomService) GetRoomById(id int64) (model.Rooms, error) {
	return r.RoomQuery.GetRoomById(id)
}

func (r *roomService) DeleteWaitList(roomId, userId int) error {
	return r.RoomQuery.DeleteWaitList(roomId, userId)
}
