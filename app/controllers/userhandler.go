package controllers

import (
	"crud/app/models"
	"crud/app/routes"
	"fmt"
	"log"

	"github.com/revel/revel"
)

type User struct {
	App
}

func (c User) Profile() revel.Result {

	if _, ok := c.Session["user"]; ok {
		username := c.Session["user"]
		user, err := models.GetUser(username)
		if err != nil {
			fmt.Println("whoops:", err)
		}
		title := user.Fullname

		return c.Render(title, user)

	}

	return c.Redirect(App.Index)

}
func (c User) ChangePassword() revel.Result {

	if _, ok := c.Session["user"]; ok {
		username := c.Session["user"]
		user, err := models.GetUser(username)
		if err != nil {
			fmt.Println("whoops:", err)
		}
		title := user.Fullname
		id := user.ID

		return c.Render(title, user, id)

	}

	return c.Redirect(App.Index)

}

/* update password */
func (c User) UpdatePassword(id int, user *models.User) revel.Result {

	user.PassValidate(c.Validation)
	// Handle errors
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		c.Flash.Error("Did not match!")
		return c.Redirect(routes.User.ChangePassword())
	}
	log.Printf("parseResponseBody password: %v\n", user.Password)
	// Ok, display the update user password

	// user.UpdatePass(c.Request)
	c.Flash.Success("Successfully Updated!")
	return c.Redirect(routes.User.Profile())

}
