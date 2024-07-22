package poll

import (
	"fmt"
	"net/http"
	"poll-app/data"
	"poll-app/utils"
	"poll-app/views/pages"

	"github.com/go-chi/chi/v5"
	"github.com/wader/gormstore/v2"
)

type PollHandler struct {
	sessionStore *gormstore.Store
}

func NewPollHandler(sessionStore *gormstore.Store) *PollHandler {
	return &PollHandler{
		sessionStore: sessionStore,
	}
}

func (h *PollHandler) InitializeRoutes(router chi.Router) {
	router.Get("/", h.PollsPage)
	router.Get("/{id}", h.PollPage)
}

func (h *PollHandler) PollsPage(w http.ResponseWriter, r *http.Request) {
	// session, _ := h.sessionStore.Get(r, "session")
	items := []data.Poll{
		{Question: "What is your favorite programming language?"},
		{Question: "Would you rather fight 1 horse-sized duck or 100 duck-sized horses?"},
	}
	err := utils.Render(pages.Polls(items), r.Context(), w)
	if err != nil {
		utils.ShowToast(w, utils.ERROR, "Failed to render page")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *PollHandler) PollPage(w http.ResponseWriter, r *http.Request) {
	pollID := chi.URLParam(r, "id")
	fmt.Println("Poll ID:", pollID)
	poll := data.Poll{
		ID:       1,
		Question: "What is your favorite programming language?",
		Options: []data.Option{
			{ID: 1, Title: "Go", TotalVotes: 10},
			{ID: 2, Title: "Python", TotalVotes: 40},
		},
		Category:  data.PollCategory{ID: 1, Name: "Programming", Label: "programming"},
		CreatedBy: data.User{ID: 1, Username: "Patrick Kabwe"},
	}
	err := utils.Render(pages.Poll(poll), r.Context(), w)
	if err != nil {
		utils.ShowToast(w, utils.ERROR, "Failed to render page")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
