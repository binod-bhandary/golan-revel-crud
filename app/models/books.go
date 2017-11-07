package models

import (
	"crud/app"
	"fmt"

	"github.com/revel/revel"
)

type Book struct {
	Isbn   string
	Title  string
	Author string
	Price  float32
}

// TODO: Make an interface for Validate() and then validation can pass in the
// key prefix ("booking.")
func (book *Book) Validate(v *revel.Validation) {

	v.Required(book.Isbn).Message("ID should be unique!")
	v.Required(book.Title).Message("Insert book title")
	v.Required(book.Author).Message("Insert book author name")
	v.Required(book.Price).Message("Insert price name")

	// v.Match(book.Price, regexp.MustCompile(`^\d+(,\d{1,2})?$`)).
	// 	Message("Credit card number must be numeric and 16 digits")

	// v.Check(book.Title,
	// 	revel.Required{},
	// 	revel.MinSize{3},
	// 	revel.MaxSize{70},
	// )
}

func AllBooks() ([]Book, error) {
	sql := "SELECT * FROM books"
	rows, err := app.DB.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bks := make([]Book, 0)
	for rows.Next() {
		bk := Book{}
		err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price) // order matters
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

func OneBook(id int) (Book, error) {
	bk := Book{}
	isbn := id

	row := app.DB.QueryRow("SELECT * FROM books WHERE isbn = $1", isbn)

	err := row.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
	if err != nil {
		return bk, err
	}

	return bk, nil
}

func (book *Book) PutBook(v *revel.Request) {
	fmt.Println(book.Isbn)
	// insert values
	app.DB.Exec("INSERT INTO books (isbn, title, author, price) VALUES ($1, $2, $3, $4)", book.Isbn, book.Title, book.Author, book.Price)
}

func (book *Book) UpdateBook(v *revel.Request) {
	// insert values
	app.DB.Exec("UPDATE books SET isbn = $1, title=$2, author=$3, price=$4 WHERE isbn=$1;", book.Isbn, book.Title, book.Author, book.Price)

}

func (book Book) DeleteBook(id int) {
	fmt.Println(id)
	app.DB.Exec("DELETE FROM books WHERE isbn=$1;", id)
}
