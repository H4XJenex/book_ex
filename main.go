package main

import (
	"book_ex/cmd/web"
	"book_ex/internal/config"
	"book_ex/internal/models"
	"crypto/tls"
	"database/sql"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"time"
)

// TODO: createReview Session Store
func main() {

	//Command-line flags
	//addr := flag.String("addr", ":8000", "HTTP network address")
	//flag.Parse()

	//Config

	cfg := config.Load()
	log.Println(cfg)

	//Logging
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(cfg)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	//Initializing the template cache
	templCache, err := web.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	//Form decoder

	formDecoder := form.NewDecoder()

	//Session manager
	sessionManager := scs.New()
	//sessionManager.Store = pqstore.New(db)
	sessionManager.Lifetime = time.Hour
	sessionManager.Cookie.Secure = true

	//
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	//Application
	app := &web.Application{
		ErrorLog:       errorLog,
		InfoLog:        infoLog,
		Reviews:        &models.ReviewModel{DB: db},
		Users:          &models.UserModel{DB: db},
		Books:          &models.BookModel{DB: db},
		TemplateCache:  templCache,
		FormDecoder:    formDecoder,
		SessionManager: sessionManager,
	}

	//Server struct
	srv := &http.Server{
		Addr:         cfg.Address,
		ErrorLog:     errorLog,
		Handler:      app.Routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  cfg.IdleTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	infoLog.Printf("Starting server on: %s", cfg.Address)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}

func openDB(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Name, cfg.ConnStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
