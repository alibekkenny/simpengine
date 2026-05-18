package romanticevent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	shared_model "github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type RomanticEventHandler struct {
	service   *RomanticEventService
	validator *validator.Validate
}

func NewRomanticEventHandler(s *RomanticEventService, v *validator.Validate) *RomanticEventHandler {
	return &RomanticEventHandler{service: s, validator: v}
}

// CreateRomanticEvent creates romantic event.
// @Summary Creates romantic event
// @Description Creates romantic event by user
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		romanticEvent  body	RomanticEventRequestDTO	true	"Romantic event data"
// @Success 	201 {object} RomanticEventResponseDTO
// @Failure 	400 {object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 {object} model.ErrorResponse "Simp target not found"
// @Failure 	500 {object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event [post]
func (h *RomanticEventHandler) CreateRomanticEvent(w http.ResponseWriter, r *http.Request) {
	var body RomanticEventRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	romanticEventID, err := h.service.CreateRomanticEvent(r.Context(), body.EventDate, body.Title, body.Description, body.SimpTargetID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", fmt.Sprintf("/romantic-event/%d", romanticEventID))
	json.NewEncoder(w).Encode(RomanticEventResponseDTO{
		ID:           romanticEventID,
		EventDate:    body.EventDate,
		Title:        body.Title,
		Description:  body.Description,
		SimpTargetID: body.SimpTargetID,
	})
}

// UpdateRomanticEvent updated romantic event by id.
// @Summary Updates romantic event
// @Description Updates an existing RomanticEvent by ID.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		romanticEvent	body	RomanticEventRequestDTO	true	"Romantic event data"
// @Param 		id  			path	int64					true	"Romantic event id"
// @Success 	204				"No Content"
// @Failure 	400 			{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 			{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 			{object} model.ErrorResponse "Simp target/Romantic event not found"
// @Failure 	409 			{object} model.ErrorResponse "Cannot edit event with given status"
// @Failure 	500 			{object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{id} [put]
func (h *RomanticEventHandler) UpdateRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body RomanticEventRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	if err := h.service.UpdateRomanticEvent(r.Context(), id, body.EventDate, body.Title, body.Description, body.SimpTargetID); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteRomanticEvent deletes RomanticEvent by ID.
// @Summary      Delete RomanticEvent
// @Description  Deletes a RomanticEvent by its ID. Only allowed if the user has proper permissions.
// @Tags         romantic_event
// @Param        id   path      int64  true  "SimpTarget ID"
// @Success      204  "No Content"
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /romantic-event/{id} [delete]
func (h *RomanticEventHandler) DeleteRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	if err := h.service.DeleteRomanticEvent(r.Context(), id); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewRomanticEvent views RomanticEvent by ID.
// @Summary      View RomanticEvent
// @Description  Views a RomanticEvent by its ID. Users can only see their romantic events.
// @Tags         romantic_event
// @Produce      json
// @Param        id   path      int64  true  "RomanticEvent ID"
// @Success      200  {object} 	RomanticEventResponseDTO	"RomanticEvent"
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Security     BearerAuth
// @Router       /romantic-event/{id} [get]
func (h *RomanticEventHandler) ViewRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	event, err := h.service.GetRomanticEventByIDAndUserID(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RomanticEventResponseDTO{
		ID:           event.ID,
		EventDate:    event.EventDate,
		Title:        event.Title,
		Description:  event.Description,
		SimpTargetID: event.SimpTargetID,
		Steps:        event.Steps,
	})
}

// ViewRomanticEventsByUser views romantic events of user.
// @Summary      View RomanticEvents by User
// @Description  Views a RomanticEvents of current User.
// @Tags         romantic_event
// @Produce 	json
// @Success		200		{array} 	model.RomanticEvent  "List of RomanticEvents"
// @Failure   	400		{object}  	model.ErrorResponse  "Invalid ID"
// @Failure   	401		{object}  	model.ErrorResponse  "Unauthorized"
// @Failure   	500 	{object}  	model.ErrorResponse  "Internal Server Error"
// @Security    BearerAuth
// @Router       /romantic-event [get]
func (h *RomanticEventHandler) ViewRomanticEventsByUser(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetRomanticEventsByUserID(r.Context())
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// PublishRomanticEvent publishes romantic event.
// @Summary Publishes romantic event
// @Description Creates public token for romantic event and changes it's status to published
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		id  path	 int64	true	"RomanticEvent Id"
// @Success 	200 {object} RomanticEventDetailResponseDTO
// @Failure 	400 {object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 {object} model.ErrorResponse "RomanticEvent not found"
// @Failure 	500 {object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{id}/publish [post]
func (h *RomanticEventHandler) PublishRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	event, choices, err := h.service.PublishRomanticEvent(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(buildDetailResponse(event, choices))
}

// AddEventStep add step to romantic event.
// @Summary 	Add event step to romantic event
// @Description Add event step to romantic event. Only Romantic event owner, can add steps to it.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  path int	true	"Romantic event id"
// @Param 		eventStep  body	EventStepsRequestDTO	true	"Event step data"
// @Success 	201 {object} EventStepsResponseDTO
// @Failure 	400 {object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 {object} model.ErrorResponse "Romantic event not found"
// @Failure 	500 {object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{event_id}/steps [post]
func (h *RomanticEventHandler) AddEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
		return
	}

	var body EventStepsRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	steps := mapDTOToSteps(body.Steps)
	createdSteps, err := h.service.AddSteps(r.Context(), steps, eventID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(EventStepsResponseDTO{
		Steps: mapStepsToDTO(createdSteps),
	})
}

// UpdateEventStep updates event step.
// @Summary 	Update event step
// @Description Updates event step of romantic event. Only Romantic event owner, can update step to it.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  	path int	true	"Romantic event id"
// @Param 		id  		path int	true	"Event step id"
// @Param 		eventStep  	body	EventStepRequestDTO	true	"Event step data"
// @Success 	204 		"No content"
// @Failure 	400 		{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 		{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 		{object} model.ErrorResponse "Romantic event/Event Step not found"
// @Failure 	500 		{object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{event_id}/steps/{id} [put]
func (h *RomanticEventHandler) UpdateEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body EventStepRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	if err := h.service.UpdateStep(r.Context(), id, body.Title, body.Description, body.StepOrder, eventID); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveEventStep removes event step.
// @Summary 	Remove event step
// @Description Removes event step of romantic event. Only Romantic event owner, can update step to it.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  	path int	true	"Romantic event id"
// @Param 		id  		path int	true	"Event step id"
// @Success 	204 		"No content"
// @Failure 	400 		{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 		{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 		{object} model.ErrorResponse "Romantic event/Event step not found"
// @Failure 	500 		{object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{event_id}/steps/{id} [delete]
func (h *RomanticEventHandler) RemoveEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	if err := h.service.RemoveStep(r.Context(), id, eventID); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddStepOption add option to event step.
// @Summary 	Add option to event step
// @Description Add option to romantic event step. Only Romantic event owner, can add option to event step.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  	path int	true	"Romantic event id"
// @Param 		step_id  	path int	true	"Event step id"
// @Param 		stepOption 	body StepOptionRequestDTO	true	"Step option data"
// @Success 	201 		{object} StepOptionResponseDTO "Step option data"
// @Failure 	400 		{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 		{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 		{object} model.ErrorResponse "Romantic event/Romantic event step not found"
// @Failure 	500 		{object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{event_id}/steps/{step_id}/options [post]
func (h *RomanticEventHandler) AddStepOption(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
		return
	}

	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid step_id", shared_model.ErrInvalidParams))
		return
	}

	var body StepOptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	id, err := h.service.AddOption(r.Context(), body.Label, body.Description, body.ImgID, eventID, stepID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StepOptionResponseDTO{
		ID:          id,
		Label:       body.Label,
		Description: body.Description,
		ImgID:       body.ImgID,
	})
}

// UpdateStepOption update event step option.
// @Summary 	Update event step option
// @Description Update option of romantic event step. Only Romantic event owner, can update event step option.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  	path int	true	"Romantic event id"
// @Param 		step_id  	path int	true	"Event step id"
// @Param 		id  		path int	true	"Id"
// @Param 		stepOption 	body StepOptionRequestDTO	true	"Step option data"
// @Success 	204 		"No content"
// @Failure 	400 		{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 		{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 		{object} model.ErrorResponse "Romantic event/Event step/Step option not found"
// @Failure 	500 		{object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{event_id}/steps/{step_id}/options/{id} [put]
func (h *RomanticEventHandler) UpdateStepOption(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
		return
	}

	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid step_id", shared_model.ErrInvalidParams))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body StepOptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	if err := h.service.UpdateOption(r.Context(), id, body.Label, body.Description, body.ImgID, eventID, stepID); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveStepOption removes event step option.
// @Summary 	Removes event step option
// @Description Remove option of romantic event step. Only Romantic event owner, can remove event step option.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  	path int	true	"Romantic event id"
// @Param 		step_id  	path int	true	"Event step id"
// @Param 		id  		path int	true	"Id"
// @Success 	204 		"No content"
// @Failure 	400 		{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 		{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 		{object} model.ErrorResponse "Romantic event/Event step/Step option not found"
// @Failure 	500 		{object} model.ErrorResponse "Internal server error"
// @Security    BearerAuth
// @Router /romantic-event/{event_id}/steps/{step_id}/options/{id} [delete]
func (h *RomanticEventHandler) RemoveStepOption(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
		return
	}

	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid step_id", shared_model.ErrInvalidParams))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	if err := h.service.RemoveOption(r.Context(), id, eventID, stepID); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewAvailableOptions views available options of user.
// @Summary      View event options by User
// @Description  Views all EventOptions of current User.
// @Tags         romantic_event
// @Produce 	json
// @Success		200		{array} 	ViewTemplateEventStepResponseDTO  "List of EventStepOptions"
// @Failure   	401		{object}  	model.ErrorResponse  "Unauthorized"
// @Failure   	500 	{object}  	model.ErrorResponse  "Internal Server Error"
// @Security    BearerAuth
// @Router       /romantic-event/steps/options [get]
func (h *RomanticEventHandler) ViewAvailableOptions(w http.ResponseWriter, r *http.Request) {
	options, err := h.service.GetAvailableOptions(r.Context())
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// AddTemplateEventStep creates a new template event step.
// @Summary      Create template event step
// @Description  Creates a new template event step (admin only). Template steps serve as base reusable steps for romantic events.
// @Tags         admin_romantic_event
// @Accept       json
// @Produce      json
// @Param        step  body     TemplateEventStepRequestDTO  true  "Template event step data"
// @Success      201   {object} TemplateEventStepResponseDTO
// @Failure      400   {object} model.ErrorResponse "Invalid request (invalid body or validation error)"
// @Failure      401   {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      403   {object} model.ErrorResponse "Forbidden (insufficient privileges)"
// @Failure      500   {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/template-event/steps [post]
func (h *RomanticEventHandler) AddTemplateEventStep(w http.ResponseWriter, r *http.Request) {
	var body TemplateEventStepRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	id, err := h.service.AddTemplateEventStep(r.Context(), body.Title, body.Description)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TemplateEventStepResponseDTO{
		ID:          id,
		Title:       body.Title,
		Description: body.Description,
	})
}

// UpdateTemplateEventStep updates an existing template event step.
// @Summary      Update template event step
// @Description  Updates an existing template event step by ID (admin only).
// @Tags         admin_romantic_event
// @Accept       json
// @Produce      json
// @Param        id    path     int64                         true  "Template event step ID"
// @Param        step  body     TemplateEventStepRequestDTO    true  "Updated template event step data"
// @Success      204   "No content"
// @Failure      400   {object} model.ErrorResponse "Invalid request (invalid ID or body)"
// @Failure      401   {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      403   {object} model.ErrorResponse "Forbidden (insufficient privileges)"
// @Failure      404   {object} model.ErrorResponse "Template event step not found"
// @Failure      500   {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/template-event/steps/{id} [put]
func (h *RomanticEventHandler) UpdateTemplateEventStep(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body TemplateEventStepRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	if err := h.service.UpdateTemplateEventStep(r.Context(), id, body.Title, body.Description, id); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AddTemplateEventStepOption adds an option to a template event step.
// @Summary      Add option to template event step
// @Description  Adds an option to an existing template event step (admin only).
// @Tags         admin_romantic_event
// @Accept       json
// @Produce      json
// @Param        id        path     int64                 true  "Template event step ID"
// @Param        option    body     StepOptionRequestDTO  true  "Step option data"
// @Success      201       {object} StepOptionResponseDTO
// @Failure      400       {object} model.ErrorResponse "Invalid request (invalid ID or body)"
// @Failure      401       {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      403       {object} model.ErrorResponse "Forbidden (insufficient privileges)"
// @Failure      404       {object} model.ErrorResponse "Template event step not found"
// @Failure      500       {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/template-event/steps/{id}/options [post]
func (h *RomanticEventHandler) AddTemplateEventStepOption(w http.ResponseWriter, r *http.Request) {
	eventStepIDStr := r.PathValue("id")
	eventStepID, err := strconv.ParseInt(eventStepIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body StepOptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	id, err := h.service.AddTemplateOption(r.Context(), body.Label, body.Description, body.ImgID, eventStepID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StepOptionResponseDTO{
		ID:          id,
		Label:       body.Label,
		Description: body.Description,
		ImgID:       body.ImgID,
	})
}

// UpdateTemplateEventStepOption updates a template event step option.
// @Summary      Update template event step option
// @Description  Updates an existing option of a template event step (admin only).
// @Tags         admin_romantic_event
// @Accept       json
// @Produce      json
// @Param        step_id  path     int64                 true  "Template event step ID"
// @Param        id       path     int64                 true  "Step option ID"
// @Param        option   body     StepOptionRequestDTO  true  "Updated option data"
// @Success      204      "No content"
// @Failure      400      {object} model.ErrorResponse "Invalid request (invalid ID or body)"
// @Failure      401      {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      403      {object} model.ErrorResponse "Forbidden (insufficient privileges)"
// @Failure      404      {object} model.ErrorResponse "Template event step or option not found"
// @Failure      500      {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /admin/template-event/steps/{step_id}/options/{id} [put]
func (h *RomanticEventHandler) UpdateTemplateEventStepOption(w http.ResponseWriter, r *http.Request) {
	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid step_id", shared_model.ErrInvalidParams))
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	var body StepOptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	if err := h.service.UpdateTemplateOption(r.Context(), id, body.Label, body.Description, body.ImgID, stepID); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewTemplateEventSteps returns all template event steps.
// @Summary      Get all template event steps
// @Description  Returns a list of all template event steps available for creating romantic events.
// @Tags         romantic_event
// @Accept       json
// @Produce      json
// @Success      200  {array}  TemplateEventStepResponseDTO
// @Failure      401  {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      500  {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /template-event/steps [get]
func (h *RomanticEventHandler) ViewTemplateEventSteps(w http.ResponseWriter, r *http.Request) {
	steps, err := h.service.GetTemplateEventSteps(r.Context())
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapTemplateStepsToDTO(steps))
}

// ViewTemplateEventStep returns a specific template event step by ID.
// @Summary      Get template event step by ID
// @Description  Returns details of a template event step with all its options.
// @Tags         romantic_event
// @Accept       json
// @Produce      json
// @Param        id  path  int64  true  "Template event step ID"
// @Success      200  {object} ViewTemplateEventStepResponseDTO
// @Failure      400  {object} model.ErrorResponse "Invalid request (invalid ID)"
// @Failure      401  {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      404  {object} model.ErrorResponse "Template event step not found"
// @Failure      500  {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /template-event/steps/{id} [get]
func (h *RomanticEventHandler) ViewTemplateEventStep(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	step, err := h.service.GetTemplateEventStep(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ViewTemplateEventStepResponseDTO{
		ID:          id,
		Title:       step.Title,
		Description: step.Description,
		Options:     mapTemplateOptionsToDTO(step.Options),
	})
}

// ViewEventSteps returns steps by event ID.
// @Summary      Get event steps by event ID
// @Description  Returns steps of event with all its options.
// @Tags         romantic_event
// @Accept       json
// @Produce      json
// @Param        event_id  path  int64  true  "Romantic event ID"
// @Success      200  {object} ViewTemplateEventStepResponseDTO
// @Failure      400  {object} model.ErrorResponse "Invalid request (invalid ID)"
// @Failure      401  {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure      404  {object} model.ErrorResponse "Template event not found"
// @Failure      500  {object} model.ErrorResponse "Internal server error"
// @Security     BearerAuth
// @Router       /romantic-event/{event_id}/steps [get]
func (h *RomanticEventHandler) ViewEventSteps(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	steps, err := h.service.GetEventSteps(r.Context(), eventID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(steps)
}

// ViewPublicRomanticEvent views RomanticEvent by public token.
// @Summary      View Published romantic event
// @Description  Views a Published romantic event by its public token. Anyone can see this romantic event.
// @Tags         public_romantic_event
// @Produce      json
// @Param        public_token   path	string  true  "RomanticEvent public token"
// @Success      200  {object} 	PublicRomanticEventResponseDTO	"Public RomanticEvent"
// @Failure      400  {object}  model.ErrorResponse  "Invalid token"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
// @Router       /public/romantic-event/{public_token} [get]
func (h *RomanticEventHandler) ViewPublicRomanticEvent(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("public_token")
	if token == "" {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid token", shared_model.ErrInvalidParams))
		return
	}

	romanticEvent, err := h.service.GetRomanticEventByPublicToken(r.Context(), token)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	var answers []EventStepChoiceResponseDTO
	if romanticEvent.Status == rmodel.StatusConfirmed {
		choices, err := h.service.GetPublicEventChoices(r.Context(), token)
		if err != nil {
			shared_model.WriteErrorResponse(w, err)
			return
		}
		answers = mapChoicesToDTO(choices)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PublicRomanticEventResponseDTO{
		ID:           romanticEvent.ID,
		Title:        romanticEvent.Title,
		Status:       romanticEvent.Status,
		Description:  romanticEvent.Description,
		EventDate:    romanticEvent.EventDate,
		PublicToken:  romanticEvent.PublicToken,
		PublishedAt:  romanticEvent.PublishedAt,
		SimpTargetID: romanticEvent.SimpTargetID,
		Steps:        romanticEvent.Steps,
		Answers:      answers,
	})
}

// SubmitPublicEventAnswers submits answers for a public RomanticEvent.
// @Summary      Submit answers for a public RomanticEvent
// @Description  Allows a simp target to submit step answer choices for a published romantic event.
// @Tags         public_romantic_event
// @Accept       json
// @Produce      json
// @Param        public_token   path      string                             true  "RomanticEvent public token"
// @Param        request        body      SubmitPublicEventAnswersRequestDTO  true  "Submitted answers"
// @Success      204            "Answers submitted successfully"
// @Failure      400            {object}  model.ErrorResponse                 "Invalid token or invalid request body"
// @Failure      404            {object}  model.ErrorResponse                 "Event not found"
// @Failure      500            {object}  model.ErrorResponse                 "Internal Server Error"
// @Router       /public/romantic-event/{public_token}/answers [post]
func (h *RomanticEventHandler) SubmitPublicEventAnswers(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("public_token")
	if token == "" {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid token", shared_model.ErrInvalidParams))
		return
	}

	var body SubmitPublicEventAnswersRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrInvalidBody, err))
		return
	}

	if err := h.validator.Struct(body); err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w:\n%v", shared_model.ErrValidation, err))
		return
	}

	choices := mapDTOToPublicEventChoices(body.EventStepAnswers)
	if err := h.service.SubmitPublicEventChoices(r.Context(), token, choices); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// AcceptPublicRomanticEvent accepts a public RomanticEvent.
// @Summary      Accept a public RomanticEvent
// @Description  Triggered when the simp target accepts the published romantic event link.
// @Tags         public_romantic_event
// @Produce      json
// @Param        public_token   path      string  true  "RomanticEvent public token"
// @Success      204            "Accepted successfully"
// @Failure      400            {object}  model.ErrorResponse  "Invalid token"
// @Failure      404            {object}  model.ErrorResponse  "Event not found"
// @Failure      500            {object}  model.ErrorResponse  "Internal Server Error"
// @Router       /public/romantic-event/{public_token}/accept [post]
func (h *RomanticEventHandler) AcceptPublicRomanticEvent(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("public_token")
	if token == "" {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid token", shared_model.ErrInvalidParams))
		return
	}

	if err := h.service.AcceptPublicRomanticEvent(r.Context(), token); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RejectPublicRomanticEvent rejects a public RomanticEvent.
// @Summary      Reject a public RomanticEvent
// @Description  Triggered when the simp target rejects the published romantic event link.
// @Tags         public_romantic_event
// @Produce      json
// @Param        public_token   path      string  true  "RomanticEvent public token"
// @Success      204            "Rejected successfully"
// @Failure      400            {object}  model.ErrorResponse  "Invalid token"
// @Failure      404            {object}  model.ErrorResponse  "Event not found"
// @Failure      500            {object}  model.ErrorResponse  "Internal Server Error"
// @Router       /public/romantic-event/{public_token}/reject [post]
func (h *RomanticEventHandler) RejectPublicRomanticEvent(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("public_token")
	if token == "" {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid token", shared_model.ErrInvalidParams))
		return
	}

	if err := h.service.RejectPublicRomanticEvent(r.Context(), token); err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RomanticEventHandler) GetEventChoices(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	choices, err := h.service.GetEventChoices(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mapChoicesToDTO(choices))
}
