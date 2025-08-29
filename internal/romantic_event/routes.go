package romanticevent

import (
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/alibekkenny/simpengine/internal/romantic_event/repository/postgres"
	"github.com/go-playground/validator/v10"
)

func RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	eventRepo := postgres.NewRomanticEventRepository(cfg.DB)
	stepRepo := postgres.NewEventStepRepository(cfg.DB)
	optionRepo := postgres.NewEventStepOptionRepository(cfg.DB)

	service := NewRomanticEventService(repo)
	validator := validator.New()
	handler := NewRomanticEventHandler(service, validator)

	mux.Handle("POST /romantic-event", auth.AuthMiddleware(http.HandlerFunc(handler.CreateRomanticEvent)))
	mux.Handle("PUT /romantic-event/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateRomanticEvent)))
	mux.Handle("DELETE /romantic-event/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DeleteRomanticEvent)))
	mux.Handle("GET /romantic-event/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.ViewRomanticEvent)))
	mux.Handle("GET /romantic-event", auth.AuthMiddleware(http.HandlerFunc(handler.ViewRomanticEventsByUser)))

	mux.Handle("POST /romantic-event/{event_id}/steps", auth.AuthMiddleware(http.HandlerFunc(handler.AddEventStep)))
	mux.Handle("PUT /romantic-event/{event_id}/steps/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateEventStep)))
	mux.Handle("DELETE /romantic-event/{event_id}/steps/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.RemoveEventStep)))

	mux.Handle("POST /romantic-event/{event_id}/steps/{step_id}/options", auth.AuthMiddleware(http.HandlerFunc(handler.AddStepOption)))
	mux.Handle("PUT /romantic-event/{event_id}/steps/{step_id}/options/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateStepOption)))
	mux.Handle("POST /romantic-event/{event_id}/steps/{step_id}/options/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.RemoveStepOption)))
}
