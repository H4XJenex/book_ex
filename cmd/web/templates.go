package web

import (
	"book_ex/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"time"
)

type TemplateData struct {
	Review          *models.Review
	Reviews         []*models.Review
	Book            *models.Book
	Books           []*models.Book
	CurrentYear     int
	Form            any
	Flash           string
	IsAuthenticated bool
	UserID          int
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func NewTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/templates/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		t, err := template.New(name).Funcs(functions).ParseFiles("./ui/templates/base.gohtml")
		if err != nil {
			return nil, err
		}

		t, err = t.ParseGlob("./ui/templates/partials/*.gohtml")
		if err != nil {
			return nil, err
		}

		t, err = t.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = t
	}

	return cache, nil
}

func (app *Application) NewTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
		//Form:        ReviewCreateForm{},
		Flash:           app.SessionManager.PopString(r.Context(), flash),
		IsAuthenticated: app.isAuthenticated(r),
		UserID:          app.SessionManager.GetInt(r.Context(), authUserID),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}
