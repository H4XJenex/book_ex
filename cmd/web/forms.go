package web

import "book_ex/internal/validator"

type ReviewCreateForm struct {
	Title               string `form:"title"`
	Text                string `form:"text"`
	BookID              int    `form:"book_id"`
	validator.Validator `form:"-"`
}

type UserSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type UserLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}
