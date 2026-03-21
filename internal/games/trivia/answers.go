package trivia

import (
	"github.com/jgoodhcg/mindmeld/internal/db"
	"github.com/jgoodhcg/mindmeld/internal/events"
	"github.com/jgoodhcg/mindmeld/internal/triviaanswer"
)

func buildAnswerDistributionFromAnswers(question db.TriviaQuestion, answers []db.GetAnswersForQuestionRow) []events.AnswerStat {
	counts := make(map[string]int, 4)
	for _, answer := range answers {
		addSelectionCount(question, counts, answer.SelectedAnswer, 1)
	}
	return answerDistribution(question, counts)
}

func buildAnswerDistributionFromStats(question db.TriviaQuestion, stats []db.GetAnswerStatsRow) []events.AnswerStat {
	counts := make(map[string]int, 4)
	for _, stat := range stats {
		addSelectionCount(question, counts, stat.SelectedAnswer, int(stat.Count))
	}
	return answerDistribution(question, counts)
}

func addSelectionCount(question db.TriviaQuestion, counts map[string]int, selection string, count int) {
	normalized := triviaanswer.NormalizeSelection(question, selection)
	if !triviaanswer.IsRecognizedSelection(question, normalized) {
		return
	}
	counts[normalized] += count
}

func answerDistribution(question db.TriviaQuestion, counts map[string]int) []events.AnswerStat {
	distribution := make([]events.AnswerStat, 0, len(counts))
	for _, option := range triviaanswer.Options(question) {
		if count, ok := counts[option.Key]; ok {
			distribution = append(distribution, events.AnswerStat{
				Answer: option.Key,
				Count:  count,
			})
		}
	}
	return distribution
}
