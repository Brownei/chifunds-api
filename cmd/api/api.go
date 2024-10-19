package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brownei/chifunds-api/service"
	"github.com/brownei/chifunds-api/service/auth"
	"github.com/brownei/chifunds-api/service/user"
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
}

func NewServer(addr string, db *sql.DB) *application {
	return &application{
		addr:         addr,
		db:           db,
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
	store := service.NewStore(a.db)
	userHandler := user.NewUserHandler(store)
	authHandler := auth.NewAuthHandler(store)

	r.Route("/v1", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			message := "ChiFunds Api"
			utils.WriteJSON(w, http.StatusOK, message)
		})

		r.Route("/users", userHandler.AllUsersRoutes)
		r.Route("/auth", authHandler.AllAuthRoutes)
	})

	log.Printf("Listening on %s", a.addr)
	return http.ListenAndServe(a.addr, r)
}
