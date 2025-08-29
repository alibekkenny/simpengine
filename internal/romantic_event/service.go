package romanticevent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alibekkenny/simpengine/internal/auth"
	rmodel "github.com/alibekkenny/simpengine/internal/romantic_event/model"
	"github.com/alibekkenny/simpengine/internal/romantic_event/repository"
	"github.com/alibekkenny/simpengine/internal/shared/model"
	"github.com/google/uuid"
)

type RomanticEventService struct {
	repo       repository.RomaticEventRepository
	stepRepo   repository.EventStepRepository
	optionRepo repository.EventStepOptionRepository
}

func NewRomanticEventService(repo repository.RomaticEventRepository, stepRepo repository.EventStepRepository, optionRepo repository.EventStepOptionRepository) *RomanticEventService {
	return &RomanticEventService{repo: repo, stepRepo: stepRepo, optionRepo: optionRepo}
}

func (s *RomanticEventService) CreateRomanticEvent(ctx context.Context, eventDate time.Time, title, description string, simpTargetID int64) (int64, error) {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return 0, model.ErrInvalidCredentials
	}

	simpTargetID, err := s.repo.CreateRomanticEvent(ctx, eventDate, title, description, simpTargetID, userID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return simpTargetID, nil
}

func (s *RomanticEventService) UpdateRomanticEvent(ctx context.Context, id int64, eventDate time.Time, title, description string, simpTargetID int64) error {
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return model.ErrInvalidCredentials
	}

	err := s.repo.UpdateRomanticEvent(ctx, id, eventDate, title, description, simpTargetID, userID)
	if err != nil {
		if errors.Is(err, model.ErrNoRecord) {
			return fmt.Errorf("%w: romantic event not found")
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

	event, err := s.repo.FindByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	steps, err := s.stepRepo.FindAllByEventID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	stepIDs := make([]int64, len(steps))
	for i, st := range steps {
		stepIDs[i] = st.ID
	}

	optionsByStep, err := s.optionRepo.FindAllByEventStepIDs(ctx, stepIDs)

	for i, step := range steps {
		steps[i].Options = optionsByStep[step.ID]
	}

	event.Steps = steps

	return event, nil
}

func (s *RomanticEventService) AddStep(ctx context.Context, title, description string, stepOrder int32, eventID int64) (int64, error) {
	_, err := s.ensureEventOwnership(ctx, eventID)
	if err != nil {
		return 0, err
	}

	id, err := s.stepRepo.CreateEventStep(ctx, title, description, stepOrder, eventID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", model.ErrInternal, err)
	}

	return id, nil
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

func (s *RomanticEventService) AddOption(ctx context.Context, label, description string, imgID uuid.UUID, eventID int64, eventStepID int64) (int64, error) {
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

	id, err := s.optionRepo.CreateEventStepOption(ctx, label, description, imgID, eventStepID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *RomanticEventService) UpdateOption(ctx context.Context, id int64, label, description string, imgID uuid.UUID, eventID int64, eventStepID int64) error {
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

	if err := s.optionRepo.UpdateEventStepOption(ctx, id, label, description, imgID); err != nil {
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
