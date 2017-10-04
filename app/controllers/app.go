package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {

	bks, err := AllBooks()
	greeting := "Aloha World"
	return c.Render(greeting)
}
