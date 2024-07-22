package api

import (
	"context"
	"log"
	"net/http"
	"poll-app/config"
	"poll-app/database"
	"poll-app/internal/auth"
	"poll-app/internal/poll"
	"poll-app/internal/user"
	"poll-app/views"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/wader/gormstore/v2"
)

type App struct {
	sessionStore *gormstore.Store
	Router       *chi.Mux
}

func (app *App) Setup() {
	// Init db
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	// Setup session store
	key := []byte("super-secret-key")

	sessionStore := gormstore.New(db, key)
	app.sessionStore = sessionStore
	// Initialize the routes
	app.Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	app.Router.Use(middleware.Logger)
	app.Router.Use(middleware.Heartbeat("/health"))
	app.Router.Use(middleware.Recoverer)
	app.Router.Use(middleware.RequestID)
	app.Router.Use(middleware.RealIP)
	app.Router.Use(middleware.AllowContentEncoding("deflate", "gzip"))
	app.Router.Use(middleware.Timeout(60))
	app.Router.Use(httprate.LimitByIP(100, 1*time.Minute))

	fs := http.FileServer(http.Dir("static"))
	app.Router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// public routes
	app.Router.Group(func(r chi.Router) {
		r.Use(app.ReverseAuthMiddleware)
		userStore := user.NewStore(db)
		authHandler := auth.NewAuthHandler(sessionStore, userStore)
		authHandler.InitializeRoutes(r)
	})

	// private routes
	app.Router.Group(func(r chi.Router) {
		r.Use(app.AuthMiddleware)
		pollHandler := poll.NewPollHandler(sessionStore)
		pollHandler.InitializeRoutes(r)
	})
}

func (app *App) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.sessionStore.Get(r, "session")
		if auth, ok := session.Values["userId"]; !ok || auth == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, views.UserKey, true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) ReverseAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.sessionStore.Get(r, "session")
		auth, ok := session.Values["userId"]
		if ok && auth != nil && r.URL.Path == "/login" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *App) Start(server *http.Server) error {
	log.Printf("Server running on port %s", config.Env.PORT)
	return server.ListenAndServe()
}
