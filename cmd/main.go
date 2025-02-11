package main

import (
	"context"

	"github.com/merynayr/AvitoShop/internal/app"
	"github.com/merynayr/AvitoShop/internal/logger"
)

// @title AvitoShop
// @version 1.0
// @description This is a AvitoShop API
// @contact.name Dmitry Boyarkin
// @contact.email boyarkin_dima2@mail.ru
func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		logger.Error("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		logger.Error("failed to run app: %s", err.Error())
	}
}
