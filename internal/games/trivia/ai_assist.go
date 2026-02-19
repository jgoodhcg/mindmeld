package trivia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
	"github.com/jgoodhcg/mindmeld/internal/questions"
)

type generatedQuestion struct {
	QuestionText  string
	CorrectAnswer string
	WrongAnswer1  string
	WrongAnswer2  string
	WrongAnswer3  string
	Source        string
}

type generateQuestionResponse struct {
	QuestionText  string `json:"question_text,omitempty"`
	CorrectAnswer string `json:"correct_answer,omitempty"`
	WrongAnswer1  string `json:"wrong_answer_1,omitempty"`
	WrongAnswer2  string `json:"wrong_answer_2,omitempty"`
	WrongAnswer3  string `json:"wrong_answer_3,omitempty"`
	Source        string `json:"source,omitempty"`
	Error         string `json:"error,omitempty"`
}

type openAIChatCompletionRequest struct {
	Model          string                 `json:"model"`
	Messages       []openAIChatMessage    `json:"messages"`
	Temperature    float64                `json:"temperature"`
	ResponseFormat map[string]interface{} `json:"response_format,omitempty"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatCompletionResponse struct {
	Choices []struct {
		Message openAIChatMessage `json:"message"`
	} `json:"choices"`
}

func generateAssistedQuestion(ctx context.Context, lobbyRating int16, topic string) (generatedQuestion, error) {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_QUESTION_ASSIST_PROVIDER")))
	key := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	model := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if model == "" {
		model = "gpt-4.1-mini"
	}

	if provider == "openai" && key != "" {
		q, err := generateQuestionWithOpenAI(ctx, key, model, lobbyRating, topic)
		if err == nil {
			return q, nil
		}
	}

	q := generateLocalQuestion(lobbyRating, topic)
	if err := validateGeneratedQuestion(q); err != nil {
		return generatedQuestion{}, err
	}
	return q, nil
}

func generateQuestionWithOpenAI(ctx context.Context, apiKey string, model string, lobbyRating int16, topic string) (generatedQuestion, error) {
	systemPrompt := strings.Join([]string{
		"You generate one multiple-choice trivia question for a social party game.",
		"Return strict JSON only with keys: question_text, correct_answer, wrong_answer_1, wrong_answer_2, wrong_answer_3.",
		"Answers must be short (max 80 chars), distinct, and plausible.",
		"Question should have one unambiguous correct answer.",
		"Keep content safe for audience: " + audiencePolicy(lobbyRating) + ".",
	}, " ")

	topic = cleanTopic(topic)
	userPrompt := "Generate a fresh trivia question."
	if topic != "" {
		userPrompt = fmt.Sprintf("Generate a fresh trivia question about: %s.", topic)
	}

	reqBody := openAIChatCompletionRequest{
		Model: model,
		Messages: []openAIChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.8,
		ResponseFormat: map[string]interface{}{
			"type": "json_object",
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return generatedQuestion{}, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, 12*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return generatedQuestion{}, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return generatedQuestion{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return generatedQuestion{}, err
	}
	if resp.StatusCode >= 400 {
		return generatedQuestion{}, fmt.Errorf("openai status %d", resp.StatusCode)
	}

	var openAIResp openAIChatCompletionResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return generatedQuestion{}, err
	}
	if len(openAIResp.Choices) == 0 {
		return generatedQuestion{}, errors.New("openai returned no choices")
	}

	var parsed struct {
		QuestionText  string `json:"question_text"`
		CorrectAnswer string `json:"correct_answer"`
		WrongAnswer1  string `json:"wrong_answer_1"`
		WrongAnswer2  string `json:"wrong_answer_2"`
		WrongAnswer3  string `json:"wrong_answer_3"`
	}
	if err := json.Unmarshal([]byte(openAIResp.Choices[0].Message.Content), &parsed); err != nil {
		return generatedQuestion{}, err
	}

	q := generatedQuestion{
		QuestionText:  strings.TrimSpace(parsed.QuestionText),
		CorrectAnswer: trimToLen(strings.TrimSpace(parsed.CorrectAnswer), 80),
		WrongAnswer1:  trimToLen(strings.TrimSpace(parsed.WrongAnswer1), 80),
		WrongAnswer2:  trimToLen(strings.TrimSpace(parsed.WrongAnswer2), 80),
		WrongAnswer3:  trimToLen(strings.TrimSpace(parsed.WrongAnswer3), 80),
		Source:        "openai",
	}

	if err := validateGeneratedQuestion(q); err != nil {
		return generatedQuestion{}, err
	}
	return q, nil
}

func generateLocalQuestion(lobbyRating int16, topic string) generatedQuestion {
	topic = cleanTopic(topic)
	candidates := localTopicCandidates(topic, lobbyRating)
	if len(candidates) == 0 {
		candidates = localTopicCandidates("", lobbyRating)
	}
	if len(candidates) == 0 {
		candidates = []generatedQuestion{
			{
				QuestionText:  "What does API stand for?",
				CorrectAnswer: "Application Programming Interface",
				WrongAnswer1:  "Automated Program Integration",
				WrongAnswer2:  "Application Process Input",
				WrongAnswer3:  "Advanced Protocol Interface",
			},
		}
	}

	index := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(candidates))
	q := candidates[index]
	q.Source = "local-fallback"
	return q
}

func localTopicCandidates(topic string, lobbyRating int16) []generatedQuestion {
	normalizedTopic := strings.ToLower(strings.TrimSpace(topic))
	allowed := questions.GetAvailableTemplates(nil, lobbyRating)

	result := make([]generatedQuestion, 0, len(allowed))
	for _, t := range allowed {
		if normalizedTopic != "" {
			candidateText := strings.ToLower(t.QuestionText + " " + t.Category)
			if !strings.Contains(candidateText, normalizedTopic) {
				continue
			}
		}
		result = append(result, generatedQuestion{
			QuestionText:  t.QuestionText,
			CorrectAnswer: trimToLen(t.CorrectAnswer, 80),
			WrongAnswer1:  trimToLen(t.WrongAnswer1, 80),
			WrongAnswer2:  trimToLen(t.WrongAnswer2, 80),
			WrongAnswer3:  trimToLen(t.WrongAnswer3, 80),
		})
	}

	return result
}

func validateGeneratedQuestion(q generatedQuestion) error {
	if strings.TrimSpace(q.QuestionText) == "" {
		return errors.New("question text is required")
	}
	answers := []string{q.CorrectAnswer, q.WrongAnswer1, q.WrongAnswer2, q.WrongAnswer3}
	seen := make(map[string]bool, len(answers))
	for _, answer := range answers {
		trimmed := strings.TrimSpace(answer)
		if trimmed == "" {
			return errors.New("all answers are required")
		}
		key := strings.ToLower(trimmed)
		if seen[key] {
			return errors.New("answers must be unique")
		}
		seen[key] = true
	}
	return nil
}

func cleanTopic(topic string) string {
	cleaned := strings.TrimSpace(topic)
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	if len(cleaned) > 60 {
		cleaned = cleaned[:60]
	}
	return cleaned
}

func trimToLen(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return strings.TrimSpace(value[:max])
}

func audiencePolicy(rating int16) string {
	switch rating {
	case contentrating.Kids:
		return "Kids: family-friendly language, no mature themes."
	case contentrating.Work:
		return "Work: professional and workplace-safe content."
	case contentrating.Adults:
		return "Adults: still avoid hateful/harassing or unsafe content."
	default:
		return "Work-safe content."
	}
}
