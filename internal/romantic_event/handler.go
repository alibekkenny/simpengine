package romanticevent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
// @Router /romantic-event [put]
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
// @Tags         simp_target
// @Produce      json
// @Param        id   path      int64  true  "RomanticEvent ID"
// @Success      200  {object} 	RomanticEventResponseDTO	"RomanticEvent"
// @Failure      400  {object}  model.ErrorResponse  "Invalid ID"
// @Failure      401  {object}  model.ErrorResponse  "Unauthorized"
// @Failure      404  {object}  model.ErrorResponse  "Not Found"
// @Failure      500  {object}  model.ErrorResponse  "Internal Server Error"
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
// @Router       /romantic-event [get]
func (h *RomanticEventHandler) ViewRomanticEventsByUser(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetRomanticEventsByUserID(r.Context())
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// PublishRomanticEvent publishes romantic event.
// @Summary Publishes romantic event
// @Description Creates public token for romantic event and changes it's status to published
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		id  path	 int64	true	"RomanticEvent Id"
// @Success 	200 {object} RomanticEventResponseDTO
// @Failure 	400 {object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 {object} model.ErrorResponse "RomanticEvent not found"
// @Failure 	500 {object} model.ErrorResponse "Internal server error"
// @Router /romantic-event/{id}/publish [post]
func (h *RomanticEventHandler) PublishRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid id", shared_model.ErrInvalidParams))
		return
	}

	status, token, err := h.service.PublishRomanticEvent(r.Context(), id)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(PublishRomanticEventResponseDTO{
		Status: status,
		Token:  token,
	})
}

// AddEventStep add step to romantic event.
// @Summary 	Add event step to romantic event
// @Description Add event step to romantic event. Only Romantic event owner, can add steps to it.
// @Tags 		romantic_event
// @Accept  	json
// @Produce  	json
// @Param 		event_id  path int	true	"Romantic event id"
// @Param 		eventStep  body	RomanticEventRequestDTO	true	"Event step data"
// @Success 	201 {object} EventStepResponseDTO
// @Failure 	400 {object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 {object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 {object} model.ErrorResponse "Romantic event not found"
// @Failure 	500 {object} model.ErrorResponse "Internal server error"
// @Router /romantic-event/{event_id}/steps [post]
func (h *RomanticEventHandler) AddEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		shared_model.WriteErrorResponse(w, fmt.Errorf("%w: invalid event_id", shared_model.ErrInvalidParams))
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

	id, err := h.service.AddStep(r.Context(), body.Title, body.Description, body.StepOrder, eventID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(EventStepResponseDTO{
		ID:          id,
		Title:       body.Title,
		Description: body.Description,
		EventOrder:  body.StepOrder,
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
// @Param 		eventStep  	body	RomanticEventRequestDTO	true	"Event step data"
// @Success 	204 		"No content"
// @Failure 	400 		{object} model.ErrorResponse "Invalid request (invalid ID, invalid body, or validation error)"
// @Failure 	401 		{object} model.ErrorResponse "Unauthorized (invalid credentials)"
// @Failure 	404 		{object} model.ErrorResponse "Romantic event/Event Step not found"
// @Failure 	500 		{object} model.ErrorResponse "Internal server error"
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

	id, err := h.service.AddOption(r.Context(), body.Label, body.ImgID, eventID, stepID)
	if err != nil {
		shared_model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(StepOptionResponseDTO{
		ID:    id,
		Label: body.Label,
		ImgID: body.ImgID,
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

	if err := h.service.UpdateOption(r.Context(), id, body.Label, body.ImgID, eventID, stepID); err != nil {
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
