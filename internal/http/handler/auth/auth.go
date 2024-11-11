package auth

import (
	"errors"
	"net/http"

	"github.com/bojackodin/notes/internal/http/encoding"
	"github.com/bojackodin/notes/internal/http/httperror"
	"github.com/bojackodin/notes/internal/log"
	"github.com/bojackodin/notes/internal/service"
)

type Controller struct {
	auth service.Auth
}

func New(auth service.Auth) *Controller {
	return &Controller{
		auth: auth,
	}
}

type signUpInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUpResponse struct {
	ID int64 `json:"id"`
}

func (ctrl *Controller) SignUp(w http.ResponseWriter, r *http.Request) error {
	logger := log.FromContext(r.Context())

	var input signUpInput
	if err := encoding.Decode(r, &input); err != nil {
		logger.Error("failed to decode body", log.Err(err))
		return httperror.WithStatusError(err, http.StatusBadRequest)
	}

	id, err := ctrl.auth.CreateUser(r.Context(), input.Username, input.Password)
	if err != nil {
		logger.Error("failed to create user", log.Err(err))
		code := http.StatusInternalServerError
		if errors.Is(err, service.ErrUserDuplicate) {
			code = http.StatusBadRequest
		}
		return httperror.WithStatusError(err, code)
	}

	_ = encoding.Encode(http.StatusCreated, w, &signUpResponse{ID: id})
	return nil
}

type signInInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

func (ctrl *Controller) SignIn(w http.ResponseWriter, r *http.Request) error {
	logger := log.FromContext(r.Context())

	var input signInInput
	if err := encoding.Decode(r, &input); err != nil {
		logger.Error("failed to decode body", log.Err(err))
		return httperror.WithStatusError(err, http.StatusBadRequest)
	}

	token, err := ctrl.auth.GenerateToken(r.Context(), input.Username, input.Password)
	if err != nil {
		logger.Error("failed to generate token", log.Err(err))
		return err
	}

	_ = encoding.Encode(http.StatusOK, w, &signInResponse{Token: token})
	return nil
}
