package api

import (
	"encoding/json"
	"go-server/internal/store"
	"go-server/internal/tokens"
	"go-server/internal/utils"
	"log"
	"net/http"
	"time"
)

type TokensHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokensHandler {
	return &TokensHandler{
		tokenStore: tokenStore,
		userStore:  userStore,
		logger:     logger,
	}
}

func (t *TokensHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		t.logger.Fatalf("HandleCreateToken error during decoding body %v", err)
		utils.WriterJSON(w, http.StatusBadRequest, utils.Envelope{"error": "error with request body, check it"})
		return
	}

	user, err := t.userStore.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		t.logger.Fatalf("HandleCreateToken cannot find user %v", err)
		utils.WriterJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	passwordDoMatches, err := user.PasswordHash.Matches(req.Password)
	if err != nil || !passwordDoMatches {
		t.logger.Fatalf("HandleCreateToken error during password matching %v", err)
		utils.WriterJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if !passwordDoMatches {
		t.logger.Fatalf("HandleCreateToken password is not matching %v", err)
		utils.WriterJSON(w, http.StatusUnauthorized, utils.Envelope{"error": "invalid credentials"})
		return
	}

	token, err := t.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		t.logger.Fatalf("HandleCreateToken cannot create token %v", err)
		utils.WriterJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriterJSON(w, http.StatusCreated, utils.Envelope{"auth_token": token})
}
