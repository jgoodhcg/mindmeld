package trivia

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
	"github.com/jgoodhcg/mindmeld/internal/questions"
)

const (
	defaultOpenAIModel     = "gpt-4.1-mini"
	defaultOpenRouterModel = "google/gemini-3.1-pro-preview"
	defaultOpenRouterTitle = "Mindmeld"
	maxAssistInputLen      = 240
	openAIRequestTimeout   = 12 * time.Second
	openRouterTimeout      = 30 * time.Second
)

var (
	openAIChatCompletionsURL     = "https://api.openai.com/v1/chat/completions"
	openRouterChatCompletionsURL = "https://openrouter.ai/api/v1/chat/completions"
	firstPersonFactPattern       = regexp.MustCompile(`(?i)^my\s+(.+?)\s+(is|are)\s+(.+)$`)
	namedFactPattern             = regexp.MustCompile(`(?i)^([a-z][a-z .'-]*?)[’']s\s+(.+?)\s+(is|are)\s+(.+)$`)
)

type statedFact struct {
	Subject   string
	Attribute string
	Verb      string
	Value     string
}

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

type chatCompletionRequest struct {
	Model          string                 `json:"model"`
	Messages       []chatMessage          `json:"messages"`
	Temperature    float64                `json:"temperature"`
	ResponseFormat map[string]interface{} `json:"response_format,omitempty"`
	Plugins        []map[string]string    `json:"plugins,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
}

type aiProviderConfig struct {
	Name      string
	APIKey    string
	Model     string
	Endpoint  string
	Headers   map[string]string
	SourceTag string
}

func generateAssistedQuestion(ctx context.Context, lobbyRating int16, topic string) (generatedQuestion, error) {
	if cfg := loadAIProviderConfig(); cfg != nil {
		q, err := generateQuestionWithProvider(ctx, *cfg, lobbyRating, topic)
		if err == nil {
			return q, nil
		}
		log.Printf("[trivia-ai] provider %s failed, falling back to local generator: %v", cfg.Name, err)
	}

	q := generateLocalQuestion(lobbyRating, topic)
	if err := validateGeneratedQuestion(q); err != nil {
		return generatedQuestion{}, err
	}
	return q, nil
}

func loadAIProviderConfig() *aiProviderConfig {
	provider := strings.ToLower(strings.TrimSpace(os.Getenv("AI_QUESTION_ASSIST_PROVIDER")))

	openRouterKey := strings.TrimSpace(os.Getenv("OPEN_ROUTER_KEY"))
	openRouterModel := strings.TrimSpace(os.Getenv("OPEN_ROUTER_MODEL"))
	if openRouterModel == "" {
		openRouterModel = defaultOpenRouterModel
	}

	openAIKey := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	openAIModel := strings.TrimSpace(os.Getenv("OPENAI_MODEL"))
	if openAIModel == "" {
		openAIModel = defaultOpenAIModel
	}

	switch provider {
	case "openrouter":
		if openRouterKey == "" {
			return nil
		}
		return &aiProviderConfig{
			Name:      "openrouter",
			APIKey:    openRouterKey,
			Model:     openRouterModel,
			Endpoint:  openRouterChatCompletionsURL,
			Headers:   openRouterHeadersFromEnv(),
			SourceTag: "openrouter",
		}
	case "openai":
		if openAIKey == "" {
			return nil
		}
		return &aiProviderConfig{
			Name:      "openai",
			APIKey:    openAIKey,
			Model:     openAIModel,
			Endpoint:  openAIChatCompletionsURL,
			SourceTag: "openai",
		}
	default:
		if openRouterKey != "" {
			return &aiProviderConfig{
				Name:      "openrouter",
				APIKey:    openRouterKey,
				Model:     openRouterModel,
				Endpoint:  openRouterChatCompletionsURL,
				Headers:   openRouterHeadersFromEnv(),
				SourceTag: "openrouter",
			}
		}
		if openAIKey != "" {
			return &aiProviderConfig{
				Name:      "openai",
				APIKey:    openAIKey,
				Model:     openAIModel,
				Endpoint:  openAIChatCompletionsURL,
				SourceTag: "openai",
			}
		}
		return nil
	}
}

func openRouterHeadersFromEnv() map[string]string {
	headers := map[string]string{
		"X-Title": defaultOpenRouterTitle,
	}

	if referer := strings.TrimSpace(os.Getenv("OPEN_ROUTER_HTTP_REFERER")); referer != "" {
		headers["HTTP-Referer"] = referer
	}
	if title := strings.TrimSpace(os.Getenv("OPEN_ROUTER_TITLE")); title != "" {
		headers["X-Title"] = title
	}

	return headers
}

func generateQuestionWithProvider(ctx context.Context, cfg aiProviderConfig, lobbyRating int16, topic string) (generatedQuestion, error) {
	systemPrompt, userPrompt := buildQuestionAssistPrompts(lobbyRating, topic)

	reqBody := chatCompletionRequest{
		Model: cfg.Model,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature:    0.8,
		ResponseFormat: questionAssistResponseFormat(),
	}
	if cfg.Name == "openrouter" {
		reqBody.Plugins = []map[string]string{
			{"id": "response-healing"},
		}
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return generatedQuestion{}, err
	}

	reqCtx, cancel := context.WithTimeout(ctx, providerTimeout(cfg.Name))
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, cfg.Endpoint, bytes.NewReader(payload))
	if err != nil {
		return generatedQuestion{}, err
	}
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range cfg.Headers {
		req.Header.Set(key, value)
	}

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
		return generatedQuestion{}, fmt.Errorf("%s status %d: %s", cfg.Name, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var completionResp chatCompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		return generatedQuestion{}, err
	}
	if len(completionResp.Choices) == 0 {
		return generatedQuestion{}, fmt.Errorf("%s returned no choices", cfg.Name)
	}

	var parsed struct {
		QuestionText  string `json:"question_text"`
		CorrectAnswer string `json:"correct_answer"`
		WrongAnswer1  string `json:"wrong_answer_1"`
		WrongAnswer2  string `json:"wrong_answer_2"`
		WrongAnswer3  string `json:"wrong_answer_3"`
	}
	if err := json.Unmarshal([]byte(completionResp.Choices[0].Message.Content), &parsed); err != nil {
		return generatedQuestion{}, err
	}

	q := generatedQuestion{
		QuestionText:  strings.TrimSpace(parsed.QuestionText),
		CorrectAnswer: trimToLen(strings.TrimSpace(parsed.CorrectAnswer), 80),
		WrongAnswer1:  trimToLen(strings.TrimSpace(parsed.WrongAnswer1), 80),
		WrongAnswer2:  trimToLen(strings.TrimSpace(parsed.WrongAnswer2), 80),
		WrongAnswer3:  trimToLen(strings.TrimSpace(parsed.WrongAnswer3), 80),
		Source:        cfg.SourceTag,
	}

	if err := validateGeneratedQuestion(q); err != nil {
		return generatedQuestion{}, err
	}
	return q, nil
}

func providerTimeout(provider string) time.Duration {
	if provider == "openrouter" {
		return openRouterTimeout
	}
	return openAIRequestTimeout
}

func questionAssistResponseFormat() map[string]interface{} {
	return map[string]interface{}{
		"type": "json_schema",
		"json_schema": map[string]interface{}{
			"name":   "trivia_question",
			"strict": true,
			"schema": map[string]interface{}{
				"type":                 "object",
				"additionalProperties": false,
				"properties": map[string]interface{}{
					"question_text": map[string]interface{}{
						"type": "string",
					},
					"correct_answer": map[string]interface{}{
						"type": "string",
					},
					"wrong_answer_1": map[string]interface{}{
						"type": "string",
					},
					"wrong_answer_2": map[string]interface{}{
						"type": "string",
					},
					"wrong_answer_3": map[string]interface{}{
						"type": "string",
					},
				},
				"required": []string{
					"question_text",
					"correct_answer",
					"wrong_answer_1",
					"wrong_answer_2",
					"wrong_answer_3",
				},
			},
		},
	}
}

func buildQuestionAssistPrompts(lobbyRating int16, topic string) (string, string) {
	systemPrompt := strings.Join([]string{
		"You generate one multiple-choice question for a social party game.",
		"Return strict JSON only with keys: question_text, correct_answer, wrong_answer_1, wrong_answer_2, wrong_answer_3.",
		"Answers must be short (max 80 chars), distinct, and plausible.",
		"Keep content safe for audience: " + audiencePolicy(lobbyRating) + ".",
		"Do not add commentary, markdown, or code fences.",
		"Output must be a single valid JSON object.",
		"Interpret the input in one of three modes.",
		"Mode 1: generic topic or category. Generate a normal trivia question about that topic with one correct answer and three plausible wrong answers.",
		"Mode 2: stated fact. Rewrite the fact into a natural player-facing question, preserve the fact exactly as the correct answer, keep any named person or subject in the question text, and generate three plausible wrong answers in the same category.",
		"If the input is an unnamed first-person fact, rewrite it using [MY_NAME] as the subject placeholder instead of generic third-person phrasing.",
		"Mode 3: personal question shell. If the input asks for a personal question but does not provide the fact, generate a clean question shell, do not invent any personal fact, and set correct_answer to exactly [fill in correct answer].",
		"Never drift into adjacent trivia when the input is clearly about a specific person or stated fact.",
		"Never fabricate personal facts.",
		"Preserve names exactly and use natural possessive grammar.",
		"If the input is noisy or misspelled, infer the likely intent conservatively without inventing facts.",
		"If multiple facts are present, choose one clear fact and make one strong question from it.",
		"Wrong answers should be similar in type, tone, and specificity to the correct answer.",
		"Prefer simple, natural wording over embellished wording.",
	}, " ")

	cleanedTopic := cleanTopic(topic)
	if cleanedTopic == "" {
		return systemPrompt, "Now generate a question from this input:\ngeneral trivia"
	}

	if fact, ok := parseStatedFact(cleanedTopic); ok && fact.Subject == "[MY_NAME]" {
		return systemPrompt, fmt.Sprintf("Now generate a question from this input:\nFirst-person stated fact about [MY_NAME]: %s", cleanedTopic)
	}

	return systemPrompt, fmt.Sprintf("Now generate a question from this input:\n%s", cleanedTopic)
}

func generateLocalQuestion(lobbyRating int16, topic string) generatedQuestion {
	topic = cleanTopic(topic)
	if fact, ok := parseStatedFact(topic); ok {
		q := generateQuestionFromStatedFact(fact)
		q.Source = "local-fallback"
		return q
	}
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
	if len(cleaned) > maxAssistInputLen {
		cleaned = cleaned[:maxAssistInputLen]
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
		return "Mild: family-friendly language, no mature themes."
	case contentrating.Work:
		return "Polite: workplace-safe and generally suitable for mixed company."
	case contentrating.Adults:
		return "Adults: still avoid hateful/harassing or unsafe content."
	default:
		return "Polite, work-safe content."
	}
}

func parseStatedFact(topic string) (statedFact, bool) {
	cleaned := strings.TrimSpace(strings.TrimRight(topic, ".!?"))
	if matches := firstPersonFactPattern.FindStringSubmatch(cleaned); len(matches) == 4 {
		return statedFact{
			Subject:   "[MY_NAME]",
			Attribute: strings.TrimSpace(matches[1]),
			Verb:      strings.ToLower(strings.TrimSpace(matches[2])),
			Value:     strings.TrimSpace(matches[3]),
		}, true
	}

	if matches := namedFactPattern.FindStringSubmatch(cleaned); len(matches) == 5 {
		return statedFact{
			Subject:   strings.TrimSpace(matches[1]),
			Attribute: strings.TrimSpace(matches[2]),
			Verb:      strings.ToLower(strings.TrimSpace(matches[3])),
			Value:     strings.TrimSpace(matches[4]),
		}, true
	}

	return statedFact{}, false
}

func generateQuestionFromStatedFact(fact statedFact) generatedQuestion {
	wrongs := distractorsForFact(fact.Attribute, fact.Value)
	return generatedQuestion{
		QuestionText:  trimToLen(buildFactQuestionText(fact), 180),
		CorrectAnswer: trimToLen(fact.Value, 80),
		WrongAnswer1:  trimToLen(wrongs[0], 80),
		WrongAnswer2:  trimToLen(wrongs[1], 80),
		WrongAnswer3:  trimToLen(wrongs[2], 80),
	}
}

func buildFactQuestionText(fact statedFact) string {
	subject := possessiveSubject(fact.Subject)
	if fact.Verb == "are" {
		return fmt.Sprintf("What are %s %s?", subject, fact.Attribute)
	}
	return fmt.Sprintf("What is %s %s?", subject, fact.Attribute)
}

func possessiveSubject(subject string) string {
	if strings.HasSuffix(subject, "s") || strings.HasSuffix(subject, "S") {
		return subject + "'"
	}
	return subject + "'s"
}

func distractorsForFact(attribute string, correct string) []string {
	pool := distractorPoolForFact(attribute, correct)
	wrongs := make([]string, 0, 3)
	seen := map[string]bool{
		strings.ToLower(strings.TrimSpace(correct)): true,
	}

	for _, candidate := range pool {
		key := strings.ToLower(strings.TrimSpace(candidate))
		if key == "" || seen[key] {
			continue
		}
		seen[key] = true
		wrongs = append(wrongs, candidate)
		if len(wrongs) == 3 {
			return wrongs
		}
	}

	fallbacks := []string{"Coffee", "Pizza", "Blue", "Dog", "Chicago", "Summer"}
	for _, candidate := range fallbacks {
		key := strings.ToLower(candidate)
		if seen[key] {
			continue
		}
		seen[key] = true
		wrongs = append(wrongs, candidate)
		if len(wrongs) == 3 {
			return wrongs
		}
	}

	for len(wrongs) < 3 {
		candidate := fmt.Sprintf("Alternative %d", len(wrongs)+1)
		key := strings.ToLower(candidate)
		if seen[key] {
			continue
		}
		seen[key] = true
		wrongs = append(wrongs, candidate)
	}

	return wrongs
}

func distractorPoolForFact(attribute string, correct string) []string {
	normalized := strings.ToLower(strings.TrimSpace(attribute + " " + correct))
	switch {
	case strings.Contains(normalized, "fruit"):
		return []string{"Apple", "Banana", "Strawberry", "Mango", "Orange"}
	case strings.Contains(normalized, "color"):
		return []string{"Blue", "Green", "Red", "Purple", "Yellow"}
	case strings.Contains(normalized, "drink"), strings.Contains(normalized, "beverage"), strings.Contains(normalized, "coffee"), strings.Contains(normalized, "tea"), strings.Contains(normalized, "soda"):
		return []string{"Coffee", "Tea", "Lemonade", "Sparkling water", "Orange juice"}
	case strings.Contains(normalized, "food"), strings.Contains(normalized, "snack"), strings.Contains(normalized, "meal"), strings.Contains(normalized, "dish"):
		return []string{"Pizza", "Tacos", "Pasta", "Sushi", "Ramen"}
	case strings.Contains(normalized, "animal"), strings.Contains(normalized, "pet"):
		return []string{"Dog", "Cat", "Rabbit", "Turtle", "Parrot"}
	case strings.Contains(normalized, "season"):
		return []string{"Spring", "Summer", "Fall", "Winter"}
	case strings.Contains(normalized, "city"), strings.Contains(normalized, "town"), strings.Contains(normalized, "hometown"):
		return []string{"Chicago", "Seattle", "Detroit", "Austin", "Boston"}
	case strings.Contains(normalized, "movie"), strings.Contains(normalized, "film"):
		return []string{"Jaws", "Alien", "Casablanca", "The Matrix", "Moonrise Kingdom"}
	case strings.Contains(normalized, "song"), strings.Contains(normalized, "karaoke"):
		return []string{"Mr. Brightside", "Dancing Queen", "Hey Jude", "Africa", "Bohemian Rhapsody"}
	case strings.Contains(normalized, "pronoun"):
		return []string{"he/him", "she/her", "xe/xem", "any pronouns", "no preference"}
	default:
		return nil
	}
}
