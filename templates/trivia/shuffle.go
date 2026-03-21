package trivia

import (
	"encoding/binary"
	"math/rand"

	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/triviaanswer"
)

// ShuffledAnswer represents an answer option with its display label
type ShuffledAnswer struct {
	Key       string // Stable answer identity
	Label     string // A, B, C, D
	Value     string // The answer text
	IsCorrect bool
}

// ShuffleAnswers returns the question's answers in a randomized order.
// The shuffle is seeded by the question ID for consistency - the same
// question always shows answers in the same shuffled order.
func ShuffleAnswers(q db.TriviaQuestion) []ShuffledAnswer {
	answers := append([]triviaanswer.Option(nil), triviaanswer.Options(q)...)

	// Seed with question ID for deterministic shuffle
	seed := int64(binary.BigEndian.Uint64(q.ID.Bytes[:8]))
	rng := rand.New(rand.NewSource(seed))

	// Fisher-Yates shuffle
	for i := len(answers) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		answers[i], answers[j] = answers[j], answers[i]
	}

	labels := []string{"A", "B", "C", "D"}
	result := make([]ShuffledAnswer, 4)
	for i, ans := range answers {
		result[i] = ShuffledAnswer{
			Key:       ans.Key,
			Label:     labels[i],
			Value:     ans.Value,
			IsCorrect: ans.IsCorrect,
		}
	}

	return result
}

func CountDistributionForAnswer(question db.TriviaQuestion, distribution []events.AnswerStat, answer ShuffledAnswer) int {
	count := 0
	for _, stat := range distribution {
		if triviaanswer.NormalizeSelection(question, stat.Answer) == answer.Key {
			count += stat.Count
		}
	}
	return count
}
