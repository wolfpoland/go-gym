package api

import (
	"encoding/json"
	"errors"
	"go-server/internal/store"
	"go-server/internal/utils"
	"log"
	"net/http"
	"regexp"
)

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (uh *UserHandler) validateRegisterRequest(registerRequest *registerUserRequest) error {
	if registerRequest.Username == "" {
		return errors.New("username is required")
	}

	if len(registerRequest.Email) > 50 {
		return errors.New("email cannot exceed 50 characters")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(registerRequest.Email) {
		return errors.New("invalid email format")
	}

	if registerRequest.Bio == "" {
		return errors.New("bio is required")
	}

	if registerRequest.Email == "" {
		return errors.New("email is required")
	}

	if registerRequest.Password == "" {
		return errors.New("passwordHash is required")
	}
	return nil
}

func (uh *UserHandler) HandleCreateUser(resWriter http.ResponseWriter, request *http.Request) {
	var requestUser = registerUserRequest{}
	err := json.NewDecoder(request.Body).Decode(&requestUser)
	if err != nil {
		uh.logger.Fatalf("Error while reading user body %v", err)
		utils.WriterJSON(resWriter, http.StatusBadRequest, utils.Envelope{"error": "Error with request body"})
		return
	}

	err = uh.validateRegisterRequest(&requestUser)
	if err != nil {
		uh.logger.Fatalf("Error validation of requestPayload failed %v", err)
		utils.WriterJSON(resWriter, http.StatusBadRequest, utils.Envelope{"error": "Error with request body validation payload failed"})
		return
	}

	user := &store.User{
		Username: requestUser.Username,
		Bio:      requestUser.Bio,
		Email:    requestUser.Email,
	}

	err = user.PasswordHash.Set(requestUser.Password)
	if err != nil {
		uh.logger.Fatalf("Error failed to hash password %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "Internal error"})
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Fatalf("Error saving user to database %v", err)
		utils.WriterJSON(resWriter, http.StatusInternalServerError, utils.Envelope{"error": "Internal error"})
		return
	}

	utils.WriterJSON(resWriter, http.StatusOK, utils.Envelope{"result": user})
}
