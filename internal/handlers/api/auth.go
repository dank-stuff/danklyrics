package api

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"codeberg.org/dankstuff/danklyrics/internal/actions"
)

type authApi struct {
	usecases *actions.Actions
}

func NewAuthApi(usecases *actions.Actions) *authApi {
	return &authApi{
		usecases: usecases,
	}
}

func (a *authApi) HandleAuth(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Email string `json:"email"`
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "invalid request body",
			SuggestedAction: "this is the request body's schema {\"email\": \"valid email string\"}",
		})
		return
	}

	if _, err := mail.ParseAddress(reqBody.Email); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "invalid email",
			SuggestedAction: "tf you think you're doing?",
		})
		return
	}

	err = a.usecases.SendVerificationEmail(reqBody.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "something went wrong",
		})
		return
	}
}

func (a *authApi) HandleConfirmAuth(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Token string `json:"token"`
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message:         "invalid request body",
			SuggestedAction: "this is the request body's schema {\"email\": \"valid email string\"}",
		})
		return
	}

	err = a.usecases.ConfirmAuth(reqBody.Token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(errorResponse{
			Message: "token is invalid or expired",
		})
		return
	}
}
