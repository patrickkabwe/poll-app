package auth

import (
	"context"
	"net/http"
	"poll-app/config"
	"poll-app/data"
	"poll-app/internal/user"
	"poll-app/utils"
	"poll-app/views"
	"poll-app/views/pages"

	"github.com/go-chi/chi/v5"
	"github.com/wader/gormstore/v2"
)

type AuthHandler struct {
	dbSessionStore *gormstore.Store
	userStore      user.UserStore
}

func NewAuthHandler(
	dbSessionStore *gormstore.Store,
	userStore user.UserStore,
) *AuthHandler {
	return &AuthHandler{
		dbSessionStore: dbSessionStore,
		userStore:      userStore,
	}
}

func (h *AuthHandler) InitializeRoutes(r chi.Router) {
	r.Get("/login", h.LoginPage)
	r.Post("/login", h.Login)
	r.Get("/register", h.RegisterPage)
	r.Post("/register", h.Register)
	r.Delete("/logout", h.Logout)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	errors := map[string]string{}
	if username == "" {
		errors["username"] = "Username is required"
	}
	if email == "" {
		errors["email"] = "Email is required"
	}
	if password == "" {
		errors["password"] = "Password is required"
	}
	if len(errors) > 0 {
		ctx := r.Context()
		ctx = context.WithValue(ctx, views.FieldValue, map[string]string{"username": username, "email": email, "password": password})
		utils.Render(pages.RegisterForm(errors), ctx, w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user := &data.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	if err := h.userStore.CreateUser(user); err != nil {
		utils.ShowToast(w, utils.ERROR, "Email already exists")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	utils.Redirect(w, r, "/login")
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	err := utils.Render(pages.Register(nil), r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	errors := map[string]string{}
	if email == "" {
		errors["email"] = "Email is required"
	}
	if password == "" {
		errors["password"] = "Password is required"
	}

	if len(errors) > 0 {
		ctx := r.Context()
		ctx = context.WithValue(ctx, views.FieldValue, map[string]string{"email": email, "password": password})
		utils.Render(pages.LoginForm(errors), ctx, w)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := h.userStore.GetByEmail(email)
	if err != nil {
		utils.ShowToast(w, utils.ERROR, "Invalid email or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !user.ComparePassword(password) {
		utils.ShowToast(w, utils.ERROR, "Invalid email or password")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	h.dbSessionStore.SessionOpts.MaxAge = 60 * 60 * 24 * 7 // 1 week
	h.dbSessionStore.SessionOpts.HttpOnly = true
	h.dbSessionStore.SessionOpts.Secure = config.Env.IsProduction()
	session, _ := h.dbSessionStore.New(r, "session")
	session.Values["userId"] = user.ID

	session.Save(r, w)
	utils.Redirect(w, r, "/")
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	utils.Render(pages.Login(nil), r.Context(), w)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := h.dbSessionStore.Get(r, "session")
	session.Options.MaxAge = -1
	session.Save(r, w)
	utils.Redirect(w, r, "/login")
}
