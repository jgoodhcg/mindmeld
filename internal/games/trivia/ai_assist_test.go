package trivia

import (
	"testing"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

func TestValidateGeneratedQuestionRequiresUniqueAnswers(t *testing.T) {
	q := generatedQuestion{
		QuestionText:  "Sample?",
		CorrectAnswer: "A",
		WrongAnswer1:  "A",
		WrongAnswer2:  "B",
		WrongAnswer3:  "C",
	}
	if err := validateGeneratedQuestion(q); err == nil {
		t.Fatal("expected duplicate answers to fail validation")
	}
}

func TestGenerateLocalQuestionReturnsValidQuestion(t *testing.T) {
	q := generateLocalQuestion(contentrating.Work, "history")
	if err := validateGeneratedQuestion(q); err != nil {
		t.Fatalf("expected generated question to be valid: %v", err)
	}
	if q.QuestionText == "" || q.CorrectAnswer == "" || q.WrongAnswer1 == "" || q.WrongAnswer2 == "" || q.WrongAnswer3 == "" {
		t.Fatal("expected all generated fields to be populated")
	}
}

func TestLocalTopicCandidatesRespectsRating(t *testing.T) {
	kids := localTopicCandidates("", contentrating.Kids)
	if len(kids) == 0 {
		t.Fatal("expected kids-safe candidates")
	}
	for _, k := range kids {
		if k.QuestionText == `In a RACI matrix, what does the "A" stand for?` {
			t.Fatal("expected work-only prompt to be excluded from kids candidates")
		}
	}
}
