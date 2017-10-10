package controllers

import (
	"crud/app/models"
	"crud/app/routes"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	*revel.Controller
}
type LoginResult struct {
	StatusCode int
	Message    string
}

var dbUsers = map[string]models.User{} // user ID, user
var dbSessions = map[string]string{}   // session ID, user ID
// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func (c App) Index() revel.Result {

	bks, err := models.AllBooks()
	if err != nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Index())
	}
	return c.Render(bks)
}

func (c App) Login() revel.Result {
	return c.Render()
}

func (c App) LoginProcess(username, password string, remember bool) revel.Result {

	// if models.alreadyLoggedIn(c.Request) {
	// 	c.Flash.Error("Already Login!")
	// 	return c.Redirect(routes.Books.Index())
	// }
	user, err := models.GetUser(username)
	//store db to variable
	dbUsers[username] = models.User{user.ID, user.Fullname, user.Email, user.Username, user.Password}
	fmt.Println(dbUsers)
	//store value to variable
	u, ok := dbUsers[username]
	if !ok {
		c.Flash.Error("Username and/or password do not match")
		return c.Redirect(routes.App.Login())
	}
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
		fmt.Println(err)
		if err == nil {
			c.Session["user"] = username
			setSession(u.Username, c.Response)
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + u.Username)
			return c.Redirect(routes.Books.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed!")
	return c.Redirect(routes.App.Login())
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

func setSession(userName string, c *revel.Response) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		c.SetCookie(cookie)
	}
}
