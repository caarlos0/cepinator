package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var db gorm.DB

func main() {
	db = NewDbConnectionPool()
	defer db.Close()
	e := echo.New()
	e.Use(mw.Logger())
	e.Use(mw.Recover())
	e.Use(mw.Gzip())
	e.Get("/ping", func(c *echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	e.Get("/ceps/:cep", func(c *echo.Context) error {
		result := db.Where("value = ?", c.Param("cep")).First(&Cep{})
		if result.RowsAffected == 0 {
			return c.JSON(http.StatusOK, nil)
		}
		return c.JSON(http.StatusOK, result.Value)
	})
	port := Getenv("PORT", "3000")
	log.Println("Running on port", port)
	e.Run(":" + port)
}
