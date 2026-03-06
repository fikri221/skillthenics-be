package auth

import (
	"nds-go-starter/internal/json"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Login godoc
// @Summary      User Login
// @Description  Authenticate user and return Access and Refresh tokens
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      loginRequest  true  "Login Credentials"
// @Success      200      {object}  LoginResponse
// @Failure      400      {object}  json.Response
// @Failure      401      {object}  json.Response
// @Router       /auth/login [post]
func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.DecodeAndValidate(r, &req); err != nil {
		json.WriteError(w, r, err)
		return
	}

	resp, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusOK, resp)
}

// Refresh godoc
// @Summary      Refresh Access Token
// @Description  Generate a new access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      refreshRequest  true  "Refresh Token Request"
// @Success      200      {object}  map[string]string "JSON with access_token"
// @Failure      400      {object}  json.Response
// @Failure      401      {object}  json.Response
// @Router       /auth/refresh [post]
func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.DecodeAndValidate(r, &req); err != nil {
		json.WriteError(w, r, err)
		return
	}

	accessToken, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		json.WriteError(w, r, err)
		return
	}

	json.Write(w, r, http.StatusOK, map[string]string{
		"access_token": accessToken,
	})
}
