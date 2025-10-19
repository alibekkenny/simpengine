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
			ID:    option.ID,
			Label: option.Label,
			ImgID: option.ImgID,
		})
	}

	return dtoOptions
}
