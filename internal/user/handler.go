package user

import (
	"encoding/json"
	"fmt"
	"net/http"

	shared_model "github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service   *UserService
	validator *validator.Validate
}

func NewUserHandler(s *UserService, v *validator.Validate) *UserHandler {
	return &UserHandler{service: s, validator: v}
}

// Register registers users.
// @Summary Register user
// @Description Registers a goddamn user
// @Tags 		user
// @Accept  	json
// @Produce  	json
// @Param 		user  	body		RegisterRequestDTO	true	"User data"
// @Success 	201 	{object} 	RegisterResponseDTO
// @Failure 	400 	{object}  	model.ErrorResponse	"Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	409 	{object} 	model.ErrorResponse "User with such email or login already exists"
// @Router /user/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequestDTO

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	ctx := r.Context()
	id, err := h.service.Register(ctx, body.Login, body.Email, body.Password)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponseDTO{
		ID:    id,
		Login: body.Login,
		Email: body.Email,
	})
}
