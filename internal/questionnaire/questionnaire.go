package questionnaire

import (
	"context"
)

type Asker interface {
	Ask(ctx context.Context) error
}

type Questionnaire struct {
	questions []Asker
}

func New(questions ...Asker) *Questionnaire {
	return &Questionnaire{
		questions: questions,
	}
}

func (i *Questionnaire) AskQuestions(ctx context.Context) error {
	for _, question := range i.questions {
		err := question.Ask(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
