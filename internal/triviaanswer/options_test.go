package triviaanswer

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jgoodhcg/mindmeld/internal/db"
)

func testQuestion(correct, wrong1, wrong2, wrong3 string) db.TriviaQuestion {
	return db.TriviaQuestion{
		ID:            pgtype.UUID{Bytes: uuid.MustParse("11111111-1111-1111-1111-111111111111"), Valid: true},
		QuestionText:  "Test?",
		CorrectAnswer: correct,
		WrongAnswer1:  wrong1,
		WrongAnswer2:  wrong2,
		WrongAnswer3:  wrong3,
	}
}

func TestNormalizeSelectionMapsUniqueLegacyAnswerText(t *testing.T) {
	question := testQuestion("Paris", "Rome", "Madrid", "Vienna")

	if got := NormalizeSelection(question, "Rome"); got != WrongAnswer1Key {
		t.Fatalf("expected unique answer text to normalize to wrong answer 1, got %q", got)
	}
}

func TestNormalizeSelectionLeavesAmbiguousDuplicateTextUnchanged(t *testing.T) {
	question := testQuestion("Blueberry", "Blueberry", "Apple", "Pear")

	if got := NormalizeSelection(question, "Blueberry"); got != "Blueberry" {
		t.Fatalf("expected ambiguous duplicate text to remain unchanged, got %q", got)
	}
}

func TestIsCorrectSelectionUsesStableKeys(t *testing.T) {
	question := testQuestion("Blueberry", "Blueberry", "Apple", "Pear")

	if !IsCorrectSelection(question, CorrectAnswerKey) {
		t.Fatal("expected correct answer key to be correct")
	}
	if IsCorrectSelection(question, WrongAnswer1Key) {
		t.Fatal("expected wrong answer key to be incorrect")
	}
}
