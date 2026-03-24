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

type CommentHandler struct {
	store *storage.CommentStore
}

func NewCommentHandler(store *storage.CommentStore) *CommentHandler {
	return &CommentHandler{
		store: store,
	}
}

func (h *CommentHandler) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// retrieve userID from context set by middleware
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "internal auth error")
		return
	}

	var newComment models.Comment

	if err := json.NewDecoder(r.Body).Decode(&newComment); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid json payload")
		return
	}
	newComment.UserID = userID
	if newComment.Comment == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	created, err := h.store.Create(r.Context(), newComment)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create comment")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}


func (h *CommentHandler) GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	comments, err := h.store.Get(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch comments")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}


func (h *CommentHandler) GetCommentsByPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	postIDInt, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "failed to parse post id")
		return
	}
	comments, err := h.store.GetByPost(r.Context(), postIDInt)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch comments by post")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)

}