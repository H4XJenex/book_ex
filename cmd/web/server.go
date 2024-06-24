package web

import (
	"book_ex/internal/models"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log"
)

type Database map[string]int

type Application struct {
	ErrorLog       *log.Logger
	InfoLog        *log.Logger
	Reviews        *models.ReviewModel
	Users          *models.UserModel
	Books          *models.BookModel
	TemplateCache  map[string]*template.Template
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}
