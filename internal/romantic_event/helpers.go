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
			ID:        step.ID,
			Title:     step.Title,
			StepOrder: step.StepOrder,
			Options:   mapOptionsToDTO(step.Options),
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
