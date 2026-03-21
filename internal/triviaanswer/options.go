package triviaanswer

import "github.com/jgoodhcg/mindmeld/internal/db"

const (
	CorrectAnswerKey = "correct_answer"
	WrongAnswer1Key  = "wrong_answer_1"
	WrongAnswer2Key  = "wrong_answer_2"
	WrongAnswer3Key  = "wrong_answer_3"
)

type Option struct {
	Key       string
	Value     string
	IsCorrect bool
}

func Options(question db.TriviaQuestion) []Option {
	return []Option{
		{Key: CorrectAnswerKey, Value: question.CorrectAnswer, IsCorrect: true},
		{Key: WrongAnswer1Key, Value: question.WrongAnswer1},
		{Key: WrongAnswer2Key, Value: question.WrongAnswer2},
		{Key: WrongAnswer3Key, Value: question.WrongAnswer3},
	}
}

func NormalizeSelection(question db.TriviaQuestion, selection string) string {
	for _, option := range Options(question) {
		if selection == option.Key {
			return option.Key
		}
	}

	matchedKey := ""
	for _, option := range Options(question) {
		if selection != option.Value {
			continue
		}
		if matchedKey != "" {
			return selection
		}
		matchedKey = option.Key
	}

	if matchedKey != "" {
		return matchedKey
	}

	return selection
}

func IsRecognizedSelection(question db.TriviaQuestion, selection string) bool {
	normalized := NormalizeSelection(question, selection)
	for _, option := range Options(question) {
		if normalized == option.Key {
			return true
		}
	}
	return false
}

func IsCorrectSelection(question db.TriviaQuestion, selection string) bool {
	return NormalizeSelection(question, selection) == CorrectAnswerKey
}
