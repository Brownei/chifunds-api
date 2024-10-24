package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brownei/chifunds-api/store"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type application struct {
	addr         string
	db           *sql.DB
	sessionStore *sessions.CookieStore
	store        store.Store
}

func NewServer(addr string, db *sql.DB, store store.Store) *application {
	return &application{
		addr:         addr,
		db:           db,
		store:        store,
		sessionStore: sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY"))),
	}
}

func (a *application) Run() error {
	r := chi.NewRouter()

	gothic.Store = a.sessionStore

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "http://localhost:8000/v1/auth/google/callback", "email", "profile"),
	)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	//All the new handlers

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			message := "ChiFunds Api"
			utils.WriteJSON(w, http.StatusOK, message)
		})

		r.Route("/users", a.AllUsersRoutes)
		r.Route("/auth", a.AllAuthRoutes)
	})

	log.Printf("Listening on %s", a.addr)
	return http.ListenAndServe(a.addr, r)
}
