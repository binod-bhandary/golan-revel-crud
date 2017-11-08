package models

import (
	"crud/app"
	"errors"
	"fmt"
	"net/http"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

var db = make(map[int]*User)

type User struct {
	ID              int
	Fullname        string
	Email           string
	Username        string
	Password        string
	PasswordConfirm string
	AccessToken     string
}

type FBUser struct {
	Name     string `json:"name"`
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GOOGLEUser struct {
	Name     string `json:"name"`
	ID       string `json:"sub"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// TODO: Make an interface for Validate() and then validation can pass in the
// key prefix ("booking.")
func (user *User) PassValidate(v *revel.Validation) {
	v.Required(user.Password).Message("Password field required.")
	v.MinSize(user.Password, 6)
	v.Required(user.PasswordConfirm).Message("Re-Password field required.")
	v.Required(user.PasswordConfirm == user.Password).Message("The passwords do not match.")
}

func AllUsers() ([]User, error) {
	rows, err := app.DB.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]User, 0)
	for rows.Next() {
		bk := User{}
		err := rows.Scan(&bk.ID, &bk.Fullname, &bk.Email, &bk.Username) // order matters
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return bks, nil
}

func AllOtherUsers(username string) ([]User, error) {

	/* select from users */
	rows, err := app.DB.Query("SELECT * FROM users where username != $1", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]User, 0)
	for rows.Next() {
		bk := User{}
		err := rows.Scan(&bk.ID, &bk.Fullname, &bk.Email, &bk.Username, &bk.Password) // order matters
		if err != nil {
			return nil, err
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	revel.INFO.Println("userList: ", bks)
	return bks, nil
}

func GetUser(username string) (User, error) {

	fmt.Println(username)
	bk := User{}
	row := app.DB.QueryRow("SELECT * FROM users WHERE username = $1 OR email = $1", username)
	err := row.Scan(&bk.ID, &bk.Fullname, &bk.Email, &bk.Username, &bk.Password)
	fmt.Println(err)
	fmt.Println(bk)
	if err != nil {
		return bk, err
	}

	return bk, nil
}

func OneUser(r *http.Request) (User, error) {

	bk := User{}
	ID := r.FormValue("id")
	if ID == "" {
		return bk, errors.New("400. Bad Request.")
	}

	row := app.DB.QueryRow("SELECT * FROM users WHERE id = $1", ID)
	err := row.Scan(&bk.ID, &bk.Fullname, &bk.Email, &bk.Username, &bk.Password)
	fmt.Println(bk)
	if err != nil {
		return bk, err
	}

	return bk, nil
}

func GetUserDetail(ID int) (User, error) {

	bk := User{}
	row := app.DB.QueryRow("SELECT * FROM users WHERE id = $1", ID)
	err := row.Scan(&bk.ID, &bk.Fullname, &bk.Email, &bk.Username, &bk.Password)
	fmt.Println(bk)
	if err != nil {
		return bk, err
	}

	return bk, nil
}

func LogUser(req *http.Request) (User, error) {

	bk := User{}
	if req.Method == http.MethodPost {
		un := req.FormValue("username")
		// is there a username?
		if un == "" {
			return bk, errors.New("400. Bad Request.")
		}

		row := app.DB.QueryRow("SELECT * FROM users WHERE username = $1", un)

		err := row.Scan(&bk.ID, &bk.Fullname, &bk.Email, &bk.Username, &bk.Password)
		fmt.Println(bk)
		if err != nil {
			return bk, err
		}

	}

	return bk, nil
}
func PutUser(r *http.Request) (User, error) {
	// get form values
	bk := User{}
	bk.Fullname = r.FormValue("fullname")
	bk.Email = r.FormValue("email")
	bk.Username = r.FormValue("username")
	// bk.Password = r.FormValue("password")
	p := r.FormValue("password")

	// validate form values
	if bk.Fullname == "" || bk.Email == "" || bk.Username == "" || p == "" {
		return bk, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	bk.Password = string(bs)
	fmt.Println(bk)
	// insert values
	_, err = app.DB.Exec("INSERT INTO users (fullname, email, username, password) VALUES ($1, $2, $3, $4)", bk.Fullname, bk.Email, bk.Username, bk.Password)

	if err != nil {
		return bk, errors.New("500. Internal Server Error." + err.Error())
	}
	return bk, nil
}

func UpdateUser(r *http.Request) (User, error) {
	// get form values
	bk := User{}
	bk.Fullname = r.FormValue("fullname")
	bk.Email = r.FormValue("email")
	bk.Username = r.FormValue("username")
	p := r.FormValue("password")
	id := r.FormValue("id")

	// validate form values
	if bk.Fullname == "" || bk.Email == "" || bk.Username == "" || p == "" {
		return bk, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	bk.Password = string(bs)

	// insert values
	_, err = app.DB.Exec("UPDATE users SET fullname = $2, email=$3, username=$4, password=$5 WHERE id=$1;", id, bk.Fullname, bk.Email, bk.Username, bk.Password)
	if err != nil {
		return bk, err
	}
	return bk, nil
}

func UpdatePass(r *http.Request) (User, error) {
	// get form values
	bk := User{}
	p := r.FormValue("password")
	id := r.FormValue("id")

	// validate form values
	if p == "" {
		return bk, errors.New("400. Bad request. All fields must be complete")
	}

	// convert form values
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	bk.Password = string(bs)

	// insert values
	_, err = app.DB.Exec("UPDATE users SET  password=$2 WHERE id=$1;", id, bk.Password)
	if err != nil {
		return bk, err
	}
	return bk, nil
}

func DeleteUser(r *http.Request) error {
	ID := r.FormValue("id")
	if ID == "" {
		return errors.New("400. Bad Request")
	}

	_, err := app.DB.Exec("DELETE FROM users WHERE id=$1;", ID)
	if err != nil {
		return errors.New("500. Internal Server Error")
	}
	return nil
}

func PutFBUser(fb FBUser) (User, error) {

	// get form values
	user := User{}
	user.Fullname = fb.Name
	user.Email = fb.Email
	user.Username = fb.ID

	// validate form values
	if user.Username == "" || user.Fullname == "" {
		return user, errors.New("400. Bad request. All fields must be complete")
	}
	p := "password"

	// convert form values
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	user.Password = string(bs)
	fmt.Println(user)
	// insert values
	_, err = app.DB.Exec("INSERT INTO users (fullname, email, username, password) VALUES ($1, $2, $3, $4)", user.Fullname, user.Email, user.Username, user.Password)

	if err != nil {
		return user, errors.New("500. Internal Server Error." + err.Error())
	}
	return user, nil
}

func PutGOOGLEUser(gg GOOGLEUser) (User, error) {

	// get form values
	user := User{}
	user.Fullname = gg.Name
	user.Email = gg.Email
	user.Username = gg.ID

	// validate form values
	if user.Username == "" || user.Fullname == "" {
		return user, errors.New("400. Bad request. All fields must be complete")
	}
	p := "password"

	// convert form values
	bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}

	user.Password = string(bs)
	fmt.Println(user)
	// insert values
	_, err = app.DB.Exec("INSERT INTO users (fullname, email, username, password) VALUES ($1, $2, $3, $4)", user.Fullname, user.Email, user.Username, user.Password)

	if err != nil {
		return user, errors.New("500. Internal Server Error." + err.Error())
	}
	return user, nil
}
