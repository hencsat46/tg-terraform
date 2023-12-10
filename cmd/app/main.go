package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"

	delivery "github.com/VanLavr/tg-bot/internal/bot/handler/http"
	"github.com/VanLavr/tg-bot/internal/bot/service"
	"github.com/VanLavr/tg-bot/internal/models"
	"github.com/VanLavr/tg-bot/pkg/env"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	url, err := env.GetTGURL("./.env")
	if err != nil {
		log.Fatal(err)
	}
	bot := new(models.Bot)
	bot.TGURL = url

	botUsecase := service.New()
	botDelivery := delivery.New(botUsecase, *bot)
	port, err := env.GetPort("./.env")
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	botDelivery.Register(e)

	go func() {
		if err := e.Start(port); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down")

}
