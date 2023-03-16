package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/with-insomnia/Hotel/internal/app"
	"github.com/with-insomnia/Hotel/internal/config"
	"github.com/with-insomnia/Hotel/internal/repository"
	"github.com/with-insomnia/Hotel/internal/service"
)

func main() {
	cfg, err := config.InitConfig("./config/config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	db, err := repository.NewDB(cfg.Database)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("cannot ping db: %v", err)
	}
	dao := repository.NewDao(db)
	authService := service.NewAuthService(dao)
	sessionService := service.NewSessionService(dao)
	postService := service.NewPostService(dao)
	userService := service.NewUserService(dao)
	roomService := service.NewRoomService(dao)
	app := app.NewAppService(authService, sessionService, userService, postService, roomService, cfg)
	server := app.Routes(cfg.Http)
	go app.ClearSession()
	go func() {
		log.Printf("server started at %s", cfg.ServerAddress)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down servers..")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown: %s\n", err)
	}
	log.Println("Server gracefully stoped")
}
