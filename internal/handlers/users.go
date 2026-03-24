package handlers

import (
	"net/http"
	"github.com/nicholaskim7/rec-programming/internal/storage"
	"github.com/nicholaskim7/rec-programming/internal/models"
	"github.com/nicholaskim7/rec-programming/internal/utils"
	"github.com/nicholaskim7/rec-programming/internal/auth"
	"github.com/nicholaskim7/rec-programming/internal/middleware"
	"encoding/json"
	"time"
)


type UserHandler struct {
	store *storage.UserStore
}

func NewUserHandler(store *storage.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var newUser models.User

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	if newUser.FirstName == "" || newUser.LastName == "" || newUser.Email == "" || newUser.Password == "" || newUser.Username == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	created, err := h.store.Create(r.Context(), newUser)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}


func (h *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.Get(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(users)
}


func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var loginPayload models.UserLoginPayload

	if err := json.NewDecoder(r.Body).Decode(&loginPayload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid json payload")
		return
	}

	if loginPayload.Username == "" || loginPayload.Password == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing required fields")
		return
	}

	user, err := h.store.Login(r.Context(), loginPayload)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to login")
		return
	}

	// generate JWT token
	token, err := auth.CreateToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "error generating session")
		return
	}
	// set http-only cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,  // javascript cannot read this (No XSS)
		Secure:   false, // set to true in production (requires https)
		Path:     "/",
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(models.UserLoginResponse{
		Token: token,
		User: user,
	})
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// create deletion cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",              // empty value
		Expires:  time.Unix(0, 0), // expired already
		MaxAge:   -1,              // tells browser to delete this cookie
		HttpOnly: true,            // javascript cannot read this (No XSS)
		Secure:   false,           // set to true in production (requires https)
		Path:     "/",
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "logged out successfully"}`))
}


func (h *UserHandler) GetMeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
        utils.RespondWithError(w, http.StatusUnauthorized, "unauthorized")
        return
    }
	user, err := h.store.GetByID(r.Context(), userID)
    if err != nil {
        utils.RespondWithError(w, http.StatusInternalServerError, "failed to fetch profile")
        return
    }
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(user)
}