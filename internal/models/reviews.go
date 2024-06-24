package models

import (
	"database/sql"
	"errors"
	"time"
)

type Review struct {
	ID        int
	Title     string
	Text      string
	Published time.Time
	BookID    int
}

type ReviewModel struct {
	DB *sql.DB
}

func (rm *ReviewModel) Get(id int) (*Review, error) {

	row := rm.DB.QueryRow(`SELECT id, title, text, published
			FROM reviews WHERE id=$1`, id)

	review := &Review{}

	err := row.Scan(&review.ID, &review.Title, &review.Text, &review.Published)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return review, nil
}

func (rm *ReviewModel) GetAll(bookID int) ([]*Review, error) {

	rows, err := rm.DB.Query(`SELECT id, title, text, published
			FROM reviews where book_id=$1`, bookID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reviews := []*Review{}

	for rows.Next() {

		review := &Review{}

		rows.Scan(&review.ID, &review.Title, &review.Text, &review.Published)
		reviews = append(reviews, review)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reviews, nil
}

// TODO
func (rm *ReviewModel) ReviewsByBookTitle(title string) (*[]Review, error) {

	return nil, nil
}

func (rm *ReviewModel) Insert(title, text string, bookID int) (int, error) {

	id := 0
	err := rm.DB.QueryRow(`insert into reviews (title, text, published, book_id)
			values ($1, $2, now(), $3) returning id`, title, text, bookID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
