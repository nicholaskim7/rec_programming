package handlers

import (
	"net/http"
	"github.com/nicholaskim7/rec-programming/internal/storage"
	"github.com/nicholaskim7/rec-programming/internal/models"
	"github.com/nicholaskim7/rec-programming/internal/utils"
	"github.com/nicholaskim7/rec-programming/internal/middleware"
	"encoding/json"
	"strconv"
)

type PostHandler struct {
	store *storage.PostStore
}

func NewPostHandler(store *storage.PostStore) *PostHandler{
	return &PostHandler{
		store: store,
	}
}

func (h *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// retrieve userID from context set by middleware
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "internal auth error")
		return
	}

	var newPost models.Post

	if err := json.NewDecoder(r.Body).Decode(&newPost); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	// force author to be the logged in user
	newPost.UserID = userID

	if newPost.Title == "" || newPost.Body == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}
	created, err := h.store.Create(r.Context(), newPost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create post")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *PostHandler) GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := h.store.Get(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch posts")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}

func (h *PostHandler) GetPostsByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	posts, err := h.store.GetByUsername(r.Context(), username)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch posts by username")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(posts)
}


func (h *PostHandler) GetPostByPostIDHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	intPostID, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to parse post id")
		return
	}

	post, err := h.store.GetByPostID(r.Context(), intPostID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch posts by username")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}