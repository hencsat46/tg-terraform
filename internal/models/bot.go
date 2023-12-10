package models

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Bot struct {
	TGURL  string
	ChatID int64
}

type BotHTTPDelivery interface {
	Register(*echo.Echo)
	HandleWebHook(echo.Context) error
	GetMe(context.Context) (*http.Response, error)
}
