package controllers

import (
	"crud/app/models"
	"crud/app/routes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
)

type App struct {
	*revel.Controller
}
type LoginResult struct {
	StatusCode int
	Message    string
}

var FACEBOOK = &oauth2.Config{

	ClientID:     "573717996154477",
	ClientSecret: "b96136347a914974c069fcbfa592fc3e",
	Scopes:       []string{"public_profile", "email", "user_friends"},
	Endpoint:     facebook.Endpoint,
	RedirectURL:  "http://localhost:9000/fbauthlogin",
}

var GOOGLE = &oauth2.Config{

	ClientID:     "970235192502-n36ijem3q6946hntrf0c8dq4jut9fsu7.apps.googleusercontent.com",
	ClientSecret: "mGqi_eUxacyVQZsN6MAYBCt4",
	Scopes: []string{
		"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
	},
	Endpoint:    google.Endpoint,
	RedirectURL: "http://localhost:9000/googleauthlogin",
}

var dbUsers = map[string]models.User{} // user ID, user
var dbSessions = map[string]string{}   // session ID, user ID

func (c App) Index() revel.Result {
	bks, err := models.AllBooks()
	if err != nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Index())
	}
	return c.Render(bks)
}

func (c App) FBLogin() string {

	authUrl := FACEBOOK.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return authUrl

}

func (c App) GOOGLELogin() string {

	authUrl := GOOGLE.AuthCodeURL("state", oauth2.AccessTypeOffline)
	return authUrl

}
func (c App) Login() revel.Result {

	c.ViewArgs["fbauthUrl"] = c.FBLogin()
	c.ViewArgs["googleauthUrl"] = c.GOOGLELogin()
	return c.Render()
}

func (c App) LoginProcess(username, password string, remember bool) revel.Result {

	// if models.alreadyLoggedIn(c.Request) {
	// 	c.Flash.Error("Already Login!")
	// 	return c.Redirect(routes.Books.Index())
	// }
	user, err := models.GetUser(username)
	// //store db to variable
	// dbUsers[username] = models.User{user.ID, user.Fullname, user.Email, user.Username, user.Password}
	// fmt.Println(dbUsers)
	// //store value to variable
	// u, ok := dbUsers[username]
	// if !ok {
	// 	c.Flash.Error("Username and/or password do not match")
	// 	return c.Redirect(routes.App.Login())
	// }
	if err == nil {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		fmt.Println(err)
		if err == nil {
			c.Session["user"] = user.Username
			if remember {
				c.Session.SetDefaultExpiration()
			} else {
				c.Session.SetNoExpiration()
			}
			c.Flash.Success("Welcome, " + user.Fullname)
			return c.Redirect(routes.Books.Index())
		}
	}

	c.Flash.Out["username"] = username
	c.Flash.Error("Login failed!")
	return c.Redirect(routes.App.Login())
}

func (c App) FBAuthLogin(code string) revel.Result {

	// We create an empty array
	fbuser := models.FBUser{}

	token, err := FACEBOOK.Exchange(oauth2.NoContext, code)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(App.Index)
	}

	revel.INFO.Println("access token")
	revel.INFO.Println(token.AccessToken)

	resp, _ := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken))

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
	}
	log.Printf("parseResponseBody: %s\n", string(response))

	err = json.Unmarshal(response, &fbuser)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	log.Printf("user detail: %v\n", fbuser)

	username := fbuser.ID
	user, err := models.GetUser(username)
	if err != nil {
		revel.INFO.Println("user not found!")
		userDetail, err := models.PutFBUser(fbuser)
		if err != nil {
			fmt.Println("user detail whoops:", err)
		}
		c.LoginProcess(userDetail.Username, "password", true)

	} else {

		c.Session["user"] = user.Username
		c.Flash.Success("Welcome, " + user.Fullname)
		return c.Redirect(routes.Books.Index())

	}

	return c.Redirect(App.Index)
}

func (c App) GOOGLEAuthLogin(code string) revel.Result {

	// We create an empty array
	googleuser := models.GOOGLEUser{}

	token, err := GOOGLE.Exchange(oauth2.NoContext, code)
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(App.Index)
	}

	client := GOOGLE.Client(oauth2.NoContext, token)
	revel.INFO.Println("access token")
	revel.INFO.Println(token.AccessToken)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		revel.ERROR.Println(err)
		return c.Redirect(App.Index)
	}

	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ReadAll: %s\n", err)
	}
	log.Printf("parseResponseBody: %s\n", string(response))

	err = json.Unmarshal(response, &googleuser)
	if err != nil {
		fmt.Println("whoops:", err)
	}
	log.Printf("user detail: %v\n", googleuser)

	username := googleuser.ID
	user, err := models.GetUser(username)
	if err != nil {
		revel.INFO.Println("user not found!")
		userDetail, err := models.PutGOOGLEUser(googleuser)
		if err != nil {
			fmt.Println("user detail whoops:", err)
		}
		c.LoginProcess(userDetail.Username, "password", true)

	} else {

		c.Session["user"] = user.Username
		c.Flash.Success("Welcome, " + user.Fullname)
		return c.Redirect(routes.Books.Index())

	}

	return c.Redirect(App.Index)
}

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.Redirect(routes.App.Index())
}

// cookie handling
func init() {

}
func (c App) EnterDemo(user, demo string) revel.Result {

	user = c.Session["user"]

	c.Validation.Required(user)
	c.Validation.Required(demo)

	if c.Validation.HasErrors() {
		c.Flash.Error("Please choose a nick name and the demonstration type.")
		return c.Redirect(App.Index)
	}

	switch demo {
	case "refresh":
		return c.Redirect("/refresh?user=%s", user)
	case "longpolling":
		return c.Redirect("/longpolling/room?user=%s", user)
	case "websocket":
		return c.Redirect("/websocket/room?user=%s", user)
	}
	return nil
}
