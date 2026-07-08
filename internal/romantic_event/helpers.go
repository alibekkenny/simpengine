package romanticevent

import "github.com/alibekkenny/simpengine/internal/romantic_event/model"

func mapTemplateStepsToDTO(steps []*model.EventStep) []*ViewTemplateEventStepResponseDTO {
	dtoSteps := []*ViewTemplateEventStepResponseDTO{}
	for _, step := range steps {
		dtoSteps = append(dtoSteps, &ViewTemplateEventStepResponseDTO{
			ID:          step.ID,
			Title:       step.Title,
			Description: step.Description,
			Options:     mapTemplateOptionsToDTO(step.Options),
		})
	}

	return dtoSteps
}

func mapTemplateOptionsToDTO(options []*model.EventStepOption) []*StepOptionResponseDTO {
	dtoOptions := []*StepOptionResponseDTO{}
	for _, option := range options {
		dtoOptions = append(dtoOptions, &StepOptionResponseDTO{
			ID:          option.ID,
			Label:       option.Label,
			Description: option.Description,
			ImgID:       option.ImgID,
		})
	}

	return dtoOptions
}

func mapDTOToSteps(dtoSteps []EventStepRequestDTO) []*model.EventStep {
	steps := []*model.EventStep{}
	for _, dtoStep := range dtoSteps {
		steps = append(steps, &model.EventStep{
			Title:       dtoStep.Title,
			Description: dtoStep.Description,
			StepOrder:   dtoStep.StepOrder,
			Options:     mapDTOToOptions(dtoStep.Options),
		})
	}

	return steps
}

func mapStepsToDTO(steps []*model.EventStep) []EventStepResponseDTO {
	dtoSteps := []EventStepResponseDTO{}
	for _, step := range steps {
		dtoSteps = append(dtoSteps, EventStepResponseDTO{
			ID:          step.ID,
			Title:       step.Title,
			Description: step.Description,
			StepOrder:   step.StepOrder,
			Options:     mapOptionsToDTO(step.Options),
		})
	}

	return dtoSteps
}

func mapDTOToOptions(dtoOptions []StepOptionRequestDTO) []*model.EventStepOption {
	options := []*model.EventStepOption{}
	for _, dtoOption := range dtoOptions {
		options = append(options, &model.EventStepOption{
			Label:       dtoOption.Label,
			Description: dtoOption.Description,
			ImgID:       dtoOption.ImgID,
		})
	}

	return options
}

func mapOptionsToDTO(options []*model.EventStepOption) []StepOptionResponseDTO {
	dtoOptions := []StepOptionResponseDTO{}
	for _, option := range options {
		dtoOptions = append(dtoOptions, StepOptionResponseDTO{
			ID:          option.ID,
			Label:       option.Label,
			Description: option.Description,
			ImgID:       option.ImgID,
		})
	}

	return dtoOptions
}

func mapDTOToPublicEventChoices(dtoChoices []SubmitStepAnswersRequestDTO) []*model.EventStepChoice {
	choices := []*model.EventStepChoice{}
	for _, dtoChoice := range dtoChoices {
		choices = append(choices, &model.EventStepChoice{
			StepID:    dtoChoice.ID,
			OptionIDs: dtoChoice.Options,
		})
	}

	return choices
}

func mapChoicesToDTO(choices []*model.EventStepChoice) []EventStepChoiceResponseDTO {
	dtoChoices := []EventStepChoiceResponseDTO{}
	for _, choice := range choices {
		dtoChoices = append(dtoChoices, EventStepChoiceResponseDTO{
			ID:        choice.ID,
			EventID:   choice.EventID,
			StepID:    choice.StepID,
			OptionIDs: choice.OptionIDs,
		})
	}

	return dtoChoices
}

func buildDetailResponse(event *model.RomanticEvent, choices []*model.EventStepChoice) RomanticEventDetailResponseDTO {
	return RomanticEventDetailResponseDTO{
		ID:           event.ID,
		Status:       event.Status,
		EventDate:    event.EventDate,
		Title:        event.Title,
		Description:  event.Description,
		PublicToken:  event.PublicToken,
		PublishedAt:  event.PublishedAt,
		SimpTargetID: event.SimpTargetID,
		Steps:        mapStepsToDTO(event.Steps),
		Answers:      mapChoicesToDTO(choices),
		RecentOpens:  []ViewSummaryDTO{},
	}
}

func mapViewSummaries(items []model.EventViewSummary) []ViewSummaryDTO {
	out := []ViewSummaryDTO{}
	for _, it := range items {
		out = append(out, ViewSummaryDTO{Device: it.Device, OS: it.OS, Browser: it.Browser, OpenedAt: it.OpenedAt})
	}
	return out
}

func buildDetailResponseWithStats(event *model.RomanticEvent, choices []*model.EventStepChoice, stats *model.EventViewStats) RomanticEventDetailResponseDTO {
	dto := buildDetailResponse(event, choices)
	if stats != nil {
		dto.Views = stats.Views
		dto.Opens = stats.Opens
		dto.LastOpenedAt = stats.LastOpenedAt
		dto.RecentOpens = mapViewSummaries(stats.RecentOpens)
	}
	return dto
}
