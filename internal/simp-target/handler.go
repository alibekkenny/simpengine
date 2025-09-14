package simptarget

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	shared_model "github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type SimpTargetHandler struct {
	service   *SimpTargetService
	validator *validator.Validate
}

func NewSimpTargetHandler(s *SimpTargetService, v *validator.Validate) *SimpTargetHandler {
	return &SimpTargetHandler{service: s, validator: v}
}

// CreateSimpTarget creates simp target.
// @Summary Creates simp target
// @Description Creates a simp target by user
// @Tags 		simp_target
// @Accept  	json
// @Produce  	json
// @Param 		simpTarget  body	SimpTargetRequestDTO	true	"Simp target data"
// @Success 	200 {object} CreateSimpTargetResponseDTO
// @Failure 	400 {object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	500 {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router /simp-target [post]
func (h *SimpTargetHandler) CreateSimpTarget(w http.ResponseWriter, r *http.Request) {
	var body SimpTargetRequestDTO
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	err = h.validator.Struct(body)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	simpTargetID, err := h.service.CreateSimpTarget(r.Context(), body.Name, body.Description)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/simp-target/%d", simpTargetID))
	json.NewEncoder(w).Encode(CreateSimpTargetResponseDTO{
		ID:          simpTargetID,
		Name:        body.Name,
		Description: body.Description,
	})
}

// UpdateSimpTarget updates an existing SimpTarget.
// @Summary      Update SimpTarget
// @Description  Updates an existing SimpTarget by ID.
// @Tags         simp_target
// @Accept       json
// @Produce      json
// @Param        id   		path      		int64                true  "SimpTarget ID"
// @Param        simpTarget body      		SimpTargetRequestDTO true  "SimpTarget data"
// @Success      204  		"No Content"
// @Failure      400  		{object} 	 	model.ErrorResponse  "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure      401  		{object}  		model.ErrorResponse  "Unauthorized"
// @Failure      404  		{object}  		model.ErrorResponse  "Not Found"
// @Failure      500  		{object}  		model.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /simp-target/{id} [put]
func (h *SimpTargetHandler) UpdateSimpTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body SimpTargetRequestDTO
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	err = h.validator.Struct(body)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	err = h.service.UpdateSimpTarget(r.Context(), id, body.Name, body.Description)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteSimpTarget deletes a SimpTarget by ID.
// @Summary      Delete SimpTarget
// @Description  Deletes a SimpTarget by its ID. Only allowed if the user has proper permissions.
// @Tags         simp_target
// @Param        id   path      int64  true  "SimpTarget ID"
// @Success      204  "No Content"
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /simp-target/{id} [delete]
func (h *SimpTargetHandler) DeleteSimpTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	err = h.service.DeleteSimpTarget(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewSimpTarget views a SimpTarget by ID.
// @Summary      View SimpTarget
// @Description  Views a SimpTarget by its ID. Users can only see their simp targets.
// @Tags         simp_target
// @Produce      json
// @Param        id   path      int64  true  "SimpTarget ID"
// @Success      200  {object} 	SimpTarget	"SimpTarget"
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /simp-target/{id} [get]
func (h *SimpTargetHandler) ViewSimpTarget(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	target, err := h.service.GetSimpTargetByIDAndUser(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(target)
}

// ViewSimpTargetByUser views a SimpTargets of user.
// @Summary      View SimpTargets by User
// @Description  Views a SimpTargets of current User.
// @Tags         simp_target
// @Produce 	json
// @Success		200		{array} 	SimpTarget "List of SimpTargets"
// @Failure   	400		{object}  	model.ErrorResponse  "Invalid ID"
// @Failure   	401		{object}  	model.ErrorResponse  "Unauthorized"
// @Failure   	500 	{object}  	model.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /simp-target [get]
func (h *SimpTargetHandler) ViewSimpTargetByUser(w http.ResponseWriter, r *http.Request) {
	targets, err := h.service.GetSimpTargetsByUserID(r.Context())
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(targets)
}
