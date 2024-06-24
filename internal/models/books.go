package models

import (
	"database/sql"
	"errors"
)

type Book struct {
	ID              int
	Title           string
	Author          string
	Description     string
	PublicationYear int
	PageQuantity    int
	Price           float64
}

type BookModel struct {
	DB *sql.DB
}

func (m *BookModel) Insert(title, author, description string, publicationYear, pageQuantity int, price float64) error {

	_, err := m.DB.Exec(`insert into books values ($1, $2, $3, $4, $5, $6)`,
		title, author, description, publicationYear, pageQuantity, price)
	if err != nil {
		return err
	}

	return nil
}

func (m *BookModel) GetByTitle(title string) (*Book, error) {

	book := &Book{}

	row := m.DB.QueryRow(`select id, title, author, description, publication_year,
       			page_quantity, price from books where title = $1`, title)

	err := row.Scan(&book.ID, &book.Author, &book.Author, &book.Description,
		&book.PublicationYear, &book.PageQuantity, &book.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return book, nil
}

func (m *BookModel) Get(id int) (*Book, error) {

	book := &Book{}

	row := m.DB.QueryRow(`select id, title, author, description, publication_year,
       			page_quantity, price from books where id = $1`, id)

	err := row.Scan(&book.ID, &book.Title, &book.Author, &book.Description,
		&book.PublicationYear, &book.PageQuantity, &book.Price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return book, nil
}

func (m *BookModel) GetAll() ([]*Book, error) {

	//TODO:Убрать ненужные поля (используются только title, author, price)
	rows, err := m.DB.Query(`select id, title, author, description, publication_year,
       			page_quantity, price from books`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	books := []*Book{}

	for rows.Next() {
		book := &Book{}

		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description, &book.PublicationYear, &book.PageQuantity, &book.Price)
		if err != nil {
			return nil, err
		}
		books = append(books, book)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}
