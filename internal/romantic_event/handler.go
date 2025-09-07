package romanticevent

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/go-playground/validator/v10"
)

type RomanticEventHandler struct {
	service   *RomanticEventService
	validator *validator.Validate
}

func NewRomanticEventHandler(s *RomanticEventService, v *validator.Validate) *RomanticEventHandler {
	return &RomanticEventHandler{service: s, validator: v}
}

// CreateRomanticEvent
func (h *RomanticEventHandler) CreateRomanticEvent(w http.ResponseWriter, r *http.Request) {
	var body RomanticEventRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err.Error()), http.StatusBadRequest)
		return
	}

	romanticEventID, err := h.service.CreateRomanticEvent(r.Context(), body.EventDate, body.Title, body.Description, body.SimpTargetID)
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/romantic-event/%d", romanticEventID))
	json.NewEncoder(w).Encode(map[string]any{
		"id":             romanticEventID,
		"event_date":     body.EventDate,
		"title":          body.Title,
		"description":    body.Description,
		"simp_target_id": body.SimpTargetID,
	})
}

// UpdateRomanticEvent
func (h *RomanticEventHandler) UpdateRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body RomanticEventRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRomanticEvent(r.Context(), id, body.EventDate, body.Title, body.Description, body.SimpTargetID); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteRomanticEvent
func (h *RomanticEventHandler) DeleteRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteRomanticEvent(r.Context(), id); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ViewRomanticEvent
func (h *RomanticEventHandler) ViewRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	event, err := h.service.GetRomanticEventByIDAndUserID(r.Context(), id)
	if err != nil {
		model.WriteErrorResponse(w, err)
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

// ViewRomanticEventByUser
func (h *RomanticEventHandler) ViewRomanticEventsByUser(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetRomanticEventsByUserID(r.Context())
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

func (h *RomanticEventHandler) PublishRomanticEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	status, token, err := h.service.PublishRomanticEvent(r.Context(), id)
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"status": status,
		"token":  token,
	})
}

func (h *RomanticEventHandler) AddEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	var body EventStepRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	id, err := h.service.AddStep(r.Context(), body.Title, body.Description, body.StepOrder, eventID)
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":          id,
		"title":       body.Title,
		"description": body.Description,
		"event_order": body.StepOrder,
	})
}

func (h *RomanticEventHandler) UpdateEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body EventStepRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateStep(r.Context(), id, body.Title, body.Description, body.StepOrder, eventID); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RomanticEventHandler) RemoveEventStep(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveStep(r.Context(), id, eventID); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RomanticEventHandler) AddStepOption(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid step_id", http.StatusBadRequest)
		return
	}

	var body StepOptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	id, err := h.service.AddOption(r.Context(), body.Label, body.ImgID, eventID, stepID)
	if err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{
		"id":     id,
		"label":  body.Label,
		"img_id": body.ImgID,
	})
}

func (h *RomanticEventHandler) UpdateStepOption(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid step_id", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body StepOptionRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, fmt.Sprintf("invalid body: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(body); err != nil {
		http.Error(w, fmt.Sprintf("validation error: %v", err), http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateOption(r.Context(), id, body.Label, body.ImgID, eventID, stepID); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RomanticEventHandler) RemoveStepOption(w http.ResponseWriter, r *http.Request) {
	eventIDStr := r.PathValue("event_id")
	eventID, err := strconv.ParseInt(eventIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid event_id", http.StatusBadRequest)
		return
	}

	stepIDStr := r.PathValue("step_id")
	stepID, err := strconv.ParseInt(stepIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid step_id", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.RemoveOption(r.Context(), id, eventID, stepID); err != nil {
		model.WriteErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
