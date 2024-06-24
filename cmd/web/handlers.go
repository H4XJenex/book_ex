package web

import (
	"book_ex/internal/models"
	"book_ex/internal/validator"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// Context keys

const (
	authUserID = "authenticatedUserID"
	flash      = "flash"
	user       = "user"
)

func (app *Application) createReview(w http.ResponseWriter, r *http.Request) {
	//TODO: Возможно есть другой способ передать ID книги
	params := httprouter.ParamsFromContext(r.Context())

	bookID, err := strconv.Atoi(params.ByName("id"))
	if err != nil || bookID < 1 {
		app.InfoLog.Println(bookID)
		app.notFound(w)
		return
	}

	data := app.NewTemplateData(r)
	data.Form = ReviewCreateForm{BookID: bookID}
	app.render(w, http.StatusOK, "create_review.gohtml", data)
}

func (app *Application) createReviewPost(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	bookID, err := strconv.Atoi(params.ByName("id"))
	if err != nil || bookID < 1 {
		app.InfoLog.Println(bookID)
		app.notFound(w)
		return
	}

	var form ReviewCreateForm

	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	//Validation check
	form.CheckFields(validator.NotBlank(form.Title), "title", "Это поле не может быть пустым")
	form.CheckFields(validator.MaxChars(form.Title, 30), "title", "Длина этого поля должны быть не больше 30 символов")
	form.CheckFields(validator.NotBlank(form.Text), "text", "Это поле не может быть пустым")

	if !form.Valid() {
		app.InfoLog.Println(form.FieldErrors)
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create_book.gohtml", data)
		return
	}

	id, err := app.Reviews.Insert(form.Title, form.Text, bookID)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Put(r.Context(), flash, "Рецензия успешно создана!")

	app.InfoLog.Printf("Id of inserted item is %d", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {

	books, err := app.Books.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.NewTemplateData(r)
	data.Books = books
	//data.Flash = app.SessionManager.PopString(r.Context(), "add")
	app.render(w, http.StatusOK, "home.gohtml", data)
}

func (app *Application) userSignup(w http.ResponseWriter, r *http.Request) {

	data := app.NewTemplateData(r)
	data.Form = UserSignupForm{}
	app.render(w, http.StatusOK, "signup.gohtml", data)
}

func (app *Application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form UserSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	app.InfoLog.Println(form.Name, form.Password, form.Email)

	form.CheckFields(validator.NotBlank(form.Name), "name", "Это поле не может быть пустым")
	form.CheckFields(validator.NotBlank(form.Email), "email", "Это поле не может быть пустым")
	form.CheckFields(validator.Matches(form.Email, validator.EmailRX), "email", "Введен некорректный e-mail")
	form.CheckFields(validator.NotBlank(form.Password), "password", "Это поле не может быть пустым")
	form.CheckFields(validator.MinChars(form.Password, 8), "password", "Пароль должен быть не меньше 8 символов")

	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.InfoLog.Println(form.FieldErrors)
		app.render(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return
	}

	err = app.Users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Этот e-mail уже используется")

			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		} else {
			app.serverError(w, err)
		}

		return
	}

	app.SessionManager.Put(r.Context(), flash, "Регистрация завершена")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *Application) userLogin(w http.ResponseWriter, r *http.Request) {

	data := app.NewTemplateData(r)
	data.Form = UserLoginForm{}
	app.render(w, http.StatusOK, "login.gohtml", data)
}

func (app *Application) userLoginPost(w http.ResponseWriter, r *http.Request) {

	var form UserLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckFields(validator.NotBlank(form.Email), "email", "Это поле не может быть пустым")
	form.CheckFields(validator.Matches(form.Email, validator.EmailRX), "email", "Некорректный почтовый адрес")
	form.CheckFields(validator.NotBlank(form.Password), "password", "Это поле не может быть пустым")

	if !form.Valid() {
		data := app.NewTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}

	id, err := app.Users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Неправильный логин или пароль")

			data := app.NewTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		} else {
			app.serverError(w, err)
		}
		return
	}

	err = app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.InfoLog.Println(id)
	app.SessionManager.Put(r.Context(), authUserID, id)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.SessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.SessionManager.Remove(r.Context(), authUserID)

	app.SessionManager.Put(r.Context(), flash, "Вы успешно вышли из аккаунта!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) createBook(w http.ResponseWriter, r *http.Request) {

}

func (app *Application) viewBook(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.InfoLog.Println(id)
		app.notFound(w)
		return
	}

	book, err := app.Books.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.ErrorLog.Print(err)
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.InfoLog.Println(id)
	reviews, err := app.Reviews.GetAll(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.ErrorLog.Print(err)
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.NewTemplateData(r)
	data.Book = book
	data.Reviews = reviews

	app.render(w, http.StatusOK, "view_book.gohtml", data)

}
