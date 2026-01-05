package romanticevent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alibekkenny/simpengine/internal/auth"
	"github.com/alibekkenny/simpengine/internal/media"
	"github.com/alibekkenny/simpengine/internal/notification"
	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/romantic_event/repository"
	"github.com/alibekkenny/simpengine/internal/shared/model"
	simptarget "github.com/alibekkenny/simpengine/internal/simp-target"
	"github.com/alibekkenny/simpengine/internal/user"
	"github.com/google/uuid"
)

type RomanticEventService struct {
	repo              repository.RomaticEventRepository
	stepRepo          repository.EventStepRepository
	optionRepo        repository.EventStepOptionRepository
	simpTargetService *simptarget.SimpTargetService
	mediaService      *media.MediaService
	userService       *user.UserService
	notifier          *notification.NotificationService
}

func NewRomanticEventService(
	repo repository.RomaticEventRepository,
	stepRepo repository.EventStepRepository,
	optionRepo repository.EventStepOptionRepository,
	simpTargetService *simptarget.SimpTargetService,
	mediaService *media.MediaService,
	userService *user.UserService,
	notifier *notification.NotificationService) *RomanticEventService {
	return &RomanticEventService{
		repo:              repo,
		stepRepo:          stepRepo,
		optionRepo:        optionRepo,
		simpTargetService: simpTargetService,
		mediaService:      mediaService,
		userService:       userService,
		notifier:          notifier,
	}
}

func (s *RomanticEventService) CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, simpTargetID int64) (int64, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return 0, model.ErrInvalidCredentials
	}

	if _, err := s.simpTargetService.GetSimpTargetByIDAndUser(ctx, simpTargetID); err != nil {
		return 0, err
	}

	simpTargetID, err := s.repo.CreateRomanticEvent(ctx, eventDate, title, description, rmodel.StatusDraft, simpTargetID, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: simp target not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return simpTargetID, nil
}

func (s *RomanticEventService) UpdateRomanticEvent(ctx context.Context, id int64, eventDate time.Time, title, description string, simpTargetID int64) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return model.ErrInvalidCredentials
	}

	if _, err := s.simpTargetService.GetSimpTargetByIDAndUser(ctx, simpTargetID); err != nil {
		return err
	}

	event, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}
	if event.Status != rmodel.StatusDraft {
		return fmt.Errorf("%w: cannot edit event with status %s", model.ErrInvalidState, event.Status)
	}

	if err := s.repo.UpdateRomanticEvent(ctx, id, eventDate, title, description, simpTargetID, userID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) DeleteRomanticEvent(ctx context.Context, id int64) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return model.ErrInvalidCredentials
	}
	// TODO: think about changing delete event steps before event itself

	err := s.repo.DeleteRomanticEvent(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) GetRomanticEventsByUserID(ctx context.Context) ([]*rmodel.RomanticEvent, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	events, err := s.repo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return events, nil
}

func (s *RomanticEventService) GetRomanticEventByIDAndUserID(ctx context.Context, id int64) (*rmodel.RomanticEvent, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	event, err := s.loadEventWithSteps(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *RomanticEventService) PublishRomanticEvent(ctx context.Context, id int64) (rmodel.RomanticEventStatus, string, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return "", "", model.ErrInvalidCredentials
	}

	event, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return "", "", fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return "", "", fmt.Errorf("%w: %v", model.ErrInternal, err)
	}
	if event.Status != rmodel.StatusDraft {
		return "", "", fmt.Errorf("%w: cannot publish event with status %s", model.ErrInvalidState, event.Status)
	}

	status := rmodel.StatusPublished
	token := uuid.New().String()

	if err := s.repo.UpdateStatusAndToken(ctx, id, userID, status, token); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return "", "", fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return "", "", fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return status, token, nil
}

func (s *RomanticEventService) AddSteps(ctx context.Context, steps []*rmodel.EventStep, eventID int64) ([]*rmodel.EventStep, error) {
	if len(steps) == 0 {
		return nil, fmt.Errorf("%w: no steps provided", model.ErrInvalidBody)
	}

	if _, err := s.ensureEventOwnership(ctx, eventID); err != nil {
		return nil, err
	}

	createdSteps, err := s.stepRepo.CreateEventStepMany(ctx, steps, eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	for i, step := range steps {
		o, err := s.optionRepo.CreateEventStepOptionMany(ctx, step.Options, createdSteps[i].ID)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
		}

		createdSteps[i].Options = o
	}

	return createdSteps, nil
}

func (s *RomanticEventService) UpdateStep(ctx context.Context, id int64, title, description string, stepOrder int32, eventID int64) error {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return err
	}

	if err := s.stepRepo.UpdateEventStep(ctx, id, title, description, stepOrder); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) RemoveStep(ctx context.Context, id int64, eventID int64) error {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return err
	}

	if err := s.stepRepo.DeleteEventStep(ctx, id); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) AddOption(ctx context.Context, label string, imgID int64, eventID int64, eventStepID int64) (int64, error) {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return 0, err
	}

	_, err = s.stepRepo.FindByIDandEventID(ctx, eventStepID, eventID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.mediaService.CheckIfExists(ctx, imgID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: image not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	id, err := s.optionRepo.CreateEventStepOption(ctx, label, imgID, eventStepID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return id, nil
}

func (s *RomanticEventService) UpdateOption(ctx context.Context, id int64, label string, imgID int64, eventID int64, eventStepID int64) error {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return err
	}

	_, err = s.stepRepo.FindByIDandEventID(ctx, eventStepID, eventID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.mediaService.CheckIfExists(ctx, imgID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: image not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.optionRepo.UpdateEventStepOption(ctx, id, label, imgID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: option not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) RemoveOption(ctx context.Context, id int64, eventID int64, eventStepID int64) error {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return err
	}

	_, err = s.stepRepo.FindByIDandEventID(ctx, eventStepID, eventID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.optionRepo.DeleteEventStepOption(ctx, id); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: option not found", model.ErrNoRecord)
		}
		return err
	}

	return nil
}

func (s *RomanticEventService) GetAvailableOptions(ctx context.Context) ([]*rmodel.EventStepOption, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, model.ErrInvalidCredentials
	}

	options, err := s.optionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return options, nil
}

func (s *RomanticEventService) AddTemplateEventStep(ctx context.Context, title, description string) (int64, error) {
	id, err := s.stepRepo.CreateTemplateEventStep(ctx, title, description)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return id, nil
}

func (s *RomanticEventService) UpdateTemplateEventStep(ctx context.Context, id int64, title, description string, eventID int64) error {
	if err := s.stepRepo.UpdateTemplateEventStep(ctx, id, title, description); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) GetTemplateEventSteps(ctx context.Context) ([]*rmodel.EventStep, error) {
	steps, err := s.stepRepo.FindAllTemplates(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.attachOptions(ctx, steps); err != nil {
		return nil, err
	}

	return steps, nil
}

func (s *RomanticEventService) AddTemplateOption(ctx context.Context, label string, imgID int64, eventStepID int64) (int64, error) {
	_, err := s.stepRepo.FindByID(ctx, eventStepID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: template event step not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.mediaService.CheckIfExists(ctx, imgID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: image not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	id, err := s.optionRepo.CreateEventStepOption(ctx, label, imgID, eventStepID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *RomanticEventService) UpdateTemplateOption(ctx context.Context, id int64, label string, imgID int64, stepID int64) error {
	_, err := s.stepRepo.FindByID(ctx, stepID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: template event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.mediaService.CheckIfExists(ctx, imgID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: image not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.optionRepo.UpdateEventStepOption(ctx, id, label, imgID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: template option not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) GetTemplateEventStep(ctx context.Context, id int64) (*rmodel.EventStep, error) {
	step, err := s.stepRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return nil, fmt.Errorf("%w: template event step not found", model.ErrNoRecord)
		}
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	options, err := s.optionRepo.FindAllByEventStepID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}
	step.Options = options

	return step, nil
}

func (s *RomanticEventService) GetEventSteps(ctx context.Context, eventID int64) ([]*rmodel.EventStep, error) {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return nil, err
	}

	steps, err := s.stepRepo.FindAllByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.attachOptions(ctx, steps); err != nil {
		return nil, err
	}

	return steps, nil
}

func (s *RomanticEventService) GetRomanticEventByPublicToken(ctx context.Context, token string) (*rmodel.RomanticEvent, error) {
	event, err := s.repo.FindByPublicToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return nil, fmt.Errorf("%w: %v", model.ErrNoRecord, err)
		}
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	steps, err := s.stepRepo.FindAllByEventID(ctx, event.ID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.attachOptions(ctx, steps); err != nil {
		return nil, err
	}
	event.Steps = steps

	return event, nil
}

func (s *RomanticEventService) SubmitPublicEventChoices(ctx context.Context, token string, choices []*rmodel.EventStepChoice) error {
	event, err := s.repo.FindByPublicToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}
	if event.Status != rmodel.StatusAccepted {
		return fmt.Errorf("%w: event is not accepted", model.ErrInternal)
	}

	for _, choice := range choices {
		existingOptions, err := s.optionRepo.FindAllByEventStepID(ctx, choice.StepID)
		if err != nil {
			return fmt.Errorf("%w: %v", model.ErrInternal, err)
		}
		if len(existingOptions) == 0 {
			return fmt.Errorf("%w: options not found", model.ErrNoRecord)
		}

		existing := make(map[int64]bool, len(existingOptions))
		for _, opt := range existingOptions {
			existing[opt.ID] = true
		}

		for _, optionID := range choice.OptionIDs {
			if _, ok := existing[optionID]; !ok {
				return fmt.Errorf("%w: option %d does not belong to this step", model.ErrInvalidParams, optionID)
			}
		}
	}

	if err := s.stepRepo.CreateAnswers(ctx, choices, event.ID); err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event step not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.repo.UpdateStatus(ctx, event.ID, event.UserID, rmodel.StatusConfirmed); err != nil {
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return nil
}

func (s *RomanticEventService) AcceptPublicRomanticEvent(ctx context.Context, token string) error {
	event, err := s.repo.FindByPublicToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}
	if event.Status != rmodel.StatusPublished {
		return fmt.Errorf("%w: event is not published", model.ErrInternal)
	}

	if err := s.repo.UpdateStatus(ctx, event.ID, event.UserID, rmodel.StatusAccepted); err != nil {
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	owner, err := s.userService.GetById(ctx, event.UserID)
	if err != nil {
		return err
	}

	if owner.NotificationsEnabled == true {
		message := event.Title + " was accepted"
		if err := s.notifier.Send(ctx, *owner, notification.ChannelTelegram, message); err != nil {
			return err
		}
	}

	return nil
}

func (s *RomanticEventService) RejectPublicRomanticEvent(ctx context.Context, token string) error {
	event, err := s.repo.FindByPublicToken(ctx, token)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: event not found", model.ErrNoRecord)
		}
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}
	if event.Status != rmodel.StatusPublished {
		return fmt.Errorf("%w: event is not published", model.ErrInternal)
	}

	if err := s.repo.UpdateStatus(ctx, event.ID, event.UserID, rmodel.StatusRejected); err != nil {
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	owner, err := s.userService.GetById(ctx, event.UserID)
	if err != nil {
		return err
	}

	if owner.NotificationsEnabled == true {
		//message := "<quote>Sometimes,\nthings don't go as planned...\nIt's called life</quote>" + event.Title + " was rejected"
		message := event.Title + " was rejected"
		if err := s.notifier.Send(ctx, *owner, notification.ChannelTelegram, message); err != nil {
			return err
		}
	}

	return nil
}

func (s *RomanticEventService) ensureEventOwnership(ctx context.Context, eventID int64) (int64, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return 0, model.ErrInvalidCredentials
	}

	_, err := s.repo.FindByIDAndUserID(ctx, eventID, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return 0, fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return userID, nil
}

func (s *RomanticEventService) loadEventWithSteps(ctx context.Context, eventID, userID int64) (*rmodel.RomanticEvent, error) {
	event, err := s.repo.FindByIDAndUserID(ctx, eventID, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return nil, fmt.Errorf("%w: romantic event not found", model.ErrNoRecord)
		}
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	steps, err := s.stepRepo.FindAllByEventID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	if err := s.attachOptions(ctx, steps); err != nil {
		return nil, err
	}

	event.Steps = steps
	return event, nil
}

func (s *RomanticEventService) attachOptions(ctx context.Context, steps []*rmodel.EventStep) error {
	if len(steps) == 0 {
		return nil
	}

	stepIDs := make([]int64, len(steps))
	for i, st := range steps {
		stepIDs[i] = st.ID
	}

	optionsByStep, err := s.optionRepo.FindAllByEventStepIDs(ctx, stepIDs)
	if err != nil {
		return fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	for i, step := range steps {
		if opts, ok := optionsByStep[step.ID]; ok {
			steps[i].Options = opts
		} else {
			steps[i].Options = []*rmodel.EventStepOption{}
		}
	}

	return nil
}

func (s *RomanticEventService) GetEventChoices(ctx context.Context, eventID int64) ([]*rmodel.EventStepChoice, error) {
	// TODO: Add authorization check if needed (e.g. check if user owns event or is the simp target)
	// _, err := s.ensureEventOwnership(ctx, eventID)
	// if err != nil {
	// 	return nil, err
	// }

	// For now, assuming handler does basic parsing and we trust the ID.
	return s.stepRepo.FindChoicesByEventID(ctx, eventID)
}
