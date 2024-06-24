package web

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"net/http"
	"runtime/debug"
)

func (app *Application) serverError(w http.ResponseWriter, err error) {

	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Output(2, trace)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *Application) render(w http.ResponseWriter, status int, page string, data *TemplateData) {

	t, ok := app.TemplateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	buf.WriteTo(w)

}

func (app *Application) decodePostForm(r *http.Request, dest any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.FormDecoder.Decode(dest, r.PostForm)
	if err != nil {
		var invalidDecErr *form.InvalidDecoderError

		if errors.As(err, &invalidDecErr) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *Application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticated).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
