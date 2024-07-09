package main

import (
	"awesomeProject/accounts"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	accountsHandler := accounts.New()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/account", accountsHandler.GetAccount)
	e.POST("/account/create", accountsHandler.CreateAccount)
	e.POST("/account/delete", accountsHandler.DeleteAccount)
	e.POST("/account/change_amount", accountsHandler.PatchAccount)
	e.POST("/account/change_name", accountsHandler.ChangeAccount)
	// Start server
	e.Logger.Fatal(e.Start(":7777"))
}
