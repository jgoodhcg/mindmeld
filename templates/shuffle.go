package templates

import (
	"encoding/binary"
	"math/rand"

	"github.com/jgoodhcg/mindmeld/internal/db"
)

// ShuffledAnswer represents an answer option with its display label
type ShuffledAnswer struct {
	Label string // A, B, C, D
	Value string // The answer text
}

// ShuffleAnswers returns the question's answers in a randomized order.
// The shuffle is seeded by the question ID for consistency - the same
// question always shows answers in the same shuffled order.
func ShuffleAnswers(q db.TriviaQuestion) []ShuffledAnswer {
	answers := []string{
		q.CorrectAnswer,
		q.WrongAnswer1,
		q.WrongAnswer2,
		q.WrongAnswer3,
	}

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
			Label: labels[i],
			Value: ans,
		}
	}

	return result
}
