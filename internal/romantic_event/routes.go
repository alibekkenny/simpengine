package romanticevent

import (
	"database/sql"
	"net/http"

	"github.com/alibekkenny/simpengine/cmd/config"
	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/alibekkenny/simpengine/internal/media"
	"github.com/alibekkenny/simpengine/internal/romantic_event/repository/postgres"
	simptarget "github.com/alibekkenny/simpengine/internal/simp-target"
	"github.com/go-playground/validator/v10"
)

type Module struct {
	Service *RomanticEventService
}

func NewModule(db *sql.DB, simpTargetService *simptarget.SimpTargetService, mediaService *media.MediaService) *Module {
	eventRepo := postgres.NewRomanticEventRepository(db)
	stepRepo := postgres.NewEventStepRepository(db)
	optionRepo := postgres.NewEventStepOptionRepository(db)

	service := NewRomanticEventService(eventRepo, stepRepo, optionRepo, simpTargetService, mediaService)

	return &Module{Service: service}
}

func (m *Module) RegisterRoutes(mux *http.ServeMux, cfg *config.Config) {
	validator := validator.New()
	handler := NewRomanticEventHandler(m.Service, validator)

	mux.Handle("POST /romantic-event", auth.AuthMiddleware(http.HandlerFunc(handler.CreateRomanticEvent)))
	mux.Handle("PUT /romantic-event/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateRomanticEvent)))
	mux.Handle("DELETE /romantic-event/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.DeleteRomanticEvent)))
	mux.Handle("GET /romantic-event/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.ViewRomanticEvent)))
	mux.Handle("GET /romantic-event", auth.AuthMiddleware(http.HandlerFunc(handler.ViewRomanticEventsByUser)))
	mux.Handle("POST /romantic-event/{id}/publish", auth.AuthMiddleware(http.HandlerFunc(handler.PublishRomanticEvent)))

	mux.Handle("GET /romantic-event/{event_id}/steps", auth.AuthMiddleware(http.HandlerFunc(handler.ViewEventSteps)))
	mux.Handle("POST /romantic-event/{event_id}/steps", auth.AuthMiddleware(http.HandlerFunc(handler.AddEventStep)))
	mux.Handle("PUT /romantic-event/{event_id}/steps/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateEventStep)))
	mux.Handle("DELETE /romantic-event/{event_id}/steps/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.RemoveEventStep)))

	mux.Handle("POST /romantic-event/{event_id}/steps/{step_id}/options", auth.AuthMiddleware(http.HandlerFunc(handler.AddStepOption)))
	mux.Handle("PUT /romantic-event/{event_id}/steps/{step_id}/options/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.UpdateStepOption)))
	mux.Handle("DELETE /romantic-event/{event_id}/steps/{step_id}/options/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.RemoveStepOption)))

	mux.Handle("GET /romantic-event/steps/options", auth.AuthMiddleware(http.HandlerFunc(handler.ViewAvailableOptions)))

	mux.Handle("POST /admin/template-event/steps",
		auth.AuthMiddleware(auth.RoleMiddleware("admin")(http.HandlerFunc(handler.AddTemplateEventStep))))
	mux.Handle("PUT /admin/template-event/steps/{id}",
		auth.AuthMiddleware(auth.RoleMiddleware("admin")(http.HandlerFunc(handler.UpdateTemplateEventStep))))
	mux.Handle("POST /admin/template-event/steps/{id}/options",
		auth.AuthMiddleware(auth.RoleMiddleware("admin")(http.HandlerFunc(handler.AddTemplateEventStepOption))))
	mux.Handle("PUT /admin/template-event/steps/{step_id}/options/{id}",
		auth.AuthMiddleware(auth.RoleMiddleware("admin")(http.HandlerFunc(handler.UpdateTemplateEventStepOption))))

	mux.Handle("GET /template-event/steps", auth.AuthMiddleware(http.HandlerFunc(handler.ViewTemplateEventSteps)))
	mux.Handle("GET /template-event/steps/{id}", auth.AuthMiddleware(http.HandlerFunc(handler.ViewTemplateEventStep)))

}
