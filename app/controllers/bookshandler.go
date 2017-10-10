package controllers

import (
	"crud/app/models"
	"crud/app/routes"
	"database/sql"
	"fmt"

	"github.com/revel/revel"
)

type Books struct {
	App
}

func (c Books) Index() revel.Result {

	Books, err := models.AllBooks()
	if err != nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.Books.Index)
	}
	return c.Render(Books)

}

func (c Books) Show(id int) revel.Result {

	book, err := models.OneBook(id)
	switch {
	case err == sql.ErrNoRows:
		c.Flash.Error("error found!")
		return c.Redirect(routes.Books.Index)
	case err != nil:
		c.Flash.Error("no data found")
		return c.NotFound("Hotel %d does not exist", id)
	}
	title := book.Title
	return c.Render(title, book)

}
func (c Books) Create() revel.Result {

	title := "create"
	return c.Render(title)
}

func (c Books) CreateProcess(book models.Book) revel.Result {

	book.Validate(c.Validation)

	// Handle errors
	fmt.Println(c.Params.Get("revise"))
	if c.Validation.HasErrors() {
		fmt.Println("failed")
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Books.Create())
	}

	// Ok, display the created user
	c.Flash.Success("%s, Successfully Created!", book.Title)
	book.PutBook(c.Request)
	return c.Redirect(routes.Books.Index())

}
func (c Books) Update(id int) revel.Result {

	book, err := models.OneBook(id)
	switch {
	case err == sql.ErrNoRows:
		c.Flash.Error("error found!")
		return c.Redirect(routes.Books.Index)
	case err != nil:
		c.Flash.Error("no data found")
		return c.NotFound("Hotel %d does not exist", id)
	}
	title := book.Title
	return c.Render(title, book, id)
}
func (c Books) UpdateProcess(id int, book models.Book) revel.Result {

	book.Validate(c.Validation)

	// Handle errors
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Books.Update(id))
	}

	// Ok, display the created user
	book.UpdateBook(c.Request)
	c.Flash.Success("%s, Successfully Updated!", book.Title)
	return c.Redirect(routes.Books.Index())

}

func (c Books) Delete(id int, book models.Book) revel.Result {

	// Ok, display the created user
	book.DeleteBook(id)
	c.Flash.Error("Successfully Delete!")
	return c.Redirect(routes.Books.Index())
}
