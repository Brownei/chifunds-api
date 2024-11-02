package api

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brownei/chifunds-api/store"
	"github.com/brownei/chifunds-api/types"
	"github.com/brownei/chifunds-api/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"go.uber.org/zap"
)

type application struct {
	addr         string
	db           *sql.DB
	sessionStore *sessions.CookieStore
	store        store.Store
	logger       *zap.SugaredLogger
	sseChannel   *types.SSEChannel
}

func NewServer(addr string, logger *zap.SugaredLogger, db *sql.DB, store store.Store) *application {
	sseChannel := &types.SSEChannel{
		Notifier: make(chan string),
		Clients:  make([]chan string, 0),
	}

	return &application{
		addr:         addr,
		db:           db,
		store:        store,
		sessionStore: sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY"))),
		logger:       logger,
		sseChannel:   sseChannel,
	}
}

func (a *application) Run() error {
	r := chi.NewRouter()
	gothic.Store = a.sessionStore

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), "https://chifunds-api.onrender.com/v1/auth/google/callback", "email", "profile"),
	)

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173", "https://chifunds.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not to preflight requests repeatedly
		Debug:            true,
	}))
	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	//All the new handlers
	r.Route("/v1", func(r chi.Router) {

		r.Get("/balance", SseRoute)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			message := "ChiFunds Api"
			utils.WriteJSON(w, http.StatusOK, message)
		})

		r.Group(func(r chi.Router) {
			r.Use(a.PublicKeyMiddleware)
			r.Get("/all-keys", a.GetAllKeysHandler)
		})

		r.Route("/users", a.AllUsersRoutes)
		r.Route("/transactions", a.AllTransactionRoutes)
		r.Route("/auth", a.AllAuthRoutes)
	})

	log.Printf("Listening on %s", a.addr)
	return http.ListenAndServe(a.addr, r)
}

func (a *application) CreateChiFundsUser() error {
	creatingNewUserQuery := `INSERT INTO "user" (email, first_name, last_name, profile_picture, password, email_verified) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, email, first_name, last_name, profile_picture, email_verified`

	_, err := a.db.Query(creatingNewUserQuery, "chifundsadmin@gmail.com", "ChiFunds", "Funding", "", "sfhkbhagassvnldfhdgklhdhguytigndnb", true)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		return err
	}

	log.Printf("Admin user created!")
	return nil
}
