package app

import (
	"github.com/with-insomnia/Hotel/internal/config"
	"github.com/with-insomnia/Hotel/internal/model"
	"github.com/with-insomnia/Hotel/internal/service"
)

var Messages model.Data

type App struct {
	authService    service.AuthService
	sessionService service.SessionService
	userService    service.UserService
	postService    service.PostService
	roomService    service.RoomService
	cfg            config.Config
}

func NewAppService(
	authService service.AuthService,
	sessionService service.SessionService,
	userService service.UserService,
	postService service.PostService,
	roomService service.RoomService,
	cfg config.Config,
) App {
	return App{
		authService:    authService,
		sessionService: sessionService,
		userService:    userService,
		postService:    postService,
		roomService:    roomService,
		cfg:            cfg,
	}
}
