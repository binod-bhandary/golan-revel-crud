package controllers

import (
	"crud/app/models"
	"crud/app/routes"
	"fmt"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	bks, err := models.AllBooks()
	fmt.Println(bks)
	if err != nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Index())
	}
	return c.Render(bks)
}
