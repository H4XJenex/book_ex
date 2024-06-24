package web

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *Application) Routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.SessionManager.LoadAndSave, app.authenticate)

	router.Handler(http.MethodGet, "/",
		dynamic.ThenFunc(app.Home))

	router.Handler(http.MethodGet, "/user/signup",
		dynamic.ThenFunc(app.userSignup))

	router.Handler(http.MethodPost, "/user/signup",
		dynamic.ThenFunc(app.userSignupPost))

	router.Handler(http.MethodGet, "/user/login",
		dynamic.ThenFunc(app.userLogin))

	router.Handler(http.MethodPost, "/user/login",
		dynamic.ThenFunc(app.userLoginPost))

	router.Handler(http.MethodGet, "/book/view/:id",
		dynamic.ThenFunc(app.viewBook))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/review/create/book/:id",
		protected.ThenFunc(app.createReview))

	router.Handler(http.MethodPost, "/review/create/book/:id",
		protected.ThenFunc(app.createReviewPost))
	router.Handler(http.MethodPost, "/user/logout",
		protected.ThenFunc(app.userLogoutPost))

	return alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then(router)
}
