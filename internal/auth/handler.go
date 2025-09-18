package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	shared_model "github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service   *AuthService
	validator *validator.Validate
}

func NewAuthHandler(s *AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{service: s, validator: v}
}

// Login authenticates users.
// @Summary User login
// @Description Authenticates a goddamn user
// @Tags 		user
// @Accept  	json
// @Produce  	json
// @Param 		user  body	LoginRequestDTO	true	"User data"
// @Success 200 {object} LoginResponseDTO
// @Failure      400  {object}  model.ErrorResponse  "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized (invalid credentials)"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Router /user/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	token, err := h.service.Login(r.Context(), body.Login, body.Password)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    token,
		Path:     "/",   // makes it available to all routes
		HttpOnly: true,  // frontend JS cannot read it (good for security)
		Secure:   false, // true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(LoginResponseDTO{
		Token: token,
	})
}
