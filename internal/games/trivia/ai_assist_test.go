package trivia

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jgoodhcg/mindmeld/internal/contentrating"
)

func TestLoadAIProviderConfigDefaultsToOpenRouterWhenKeyPresent(t *testing.T) {
	t.Setenv("OPEN_ROUTER_KEY", "test-key")
	t.Setenv("AI_QUESTION_ASSIST_PROVIDER", "")
	t.Setenv("OPEN_ROUTER_MODEL", "")
	t.Setenv("OPEN_ROUTER_TITLE", "")
	t.Setenv("OPEN_ROUTER_HTTP_REFERER", "")
	t.Setenv("OPENAI_API_KEY", "")

	cfg := loadAIProviderConfig()
	if cfg == nil {
		t.Fatal("expected provider config")
	}
	if cfg.Name != "openrouter" {
		t.Fatalf("expected openrouter provider, got %q", cfg.Name)
	}
	if cfg.Model != defaultOpenRouterModel {
		t.Fatalf("expected default model %q, got %q", defaultOpenRouterModel, cfg.Model)
	}
	if cfg.Headers["X-Title"] != defaultOpenRouterTitle {
		t.Fatalf("expected default title %q, got %q", defaultOpenRouterTitle, cfg.Headers["X-Title"])
	}
}

func TestGenerateQuestionWithProviderUsesOpenRouterCompatibleRequest(t *testing.T) {
	var seenAuthorization string
	var seenReferer string
	var seenTitle string
	var seenModel string
	var seenSystemPrompt string
	var seenUserPrompt string
	var seenResponseFormat map[string]interface{}
	var seenPlugins []map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuthorization = r.Header.Get("Authorization")
		seenReferer = r.Header.Get("HTTP-Referer")
		seenTitle = r.Header.Get("X-Title")

		var req chatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		seenModel = req.Model
		if len(req.Messages) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(req.Messages))
		}
		seenSystemPrompt = req.Messages[0].Content
		seenUserPrompt = req.Messages[1].Content
		seenResponseFormat = req.ResponseFormat
		seenPlugins = req.Plugins

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(chatCompletionResponse{
			Choices: []struct {
				Message chatMessage `json:"message"`
			}{
				{
					Message: chatMessage{
						Role:    "assistant",
						Content: `{"question_text":"What is the capital of France?","correct_answer":"Paris","wrong_answer_1":"Lyon","wrong_answer_2":"Marseille","wrong_answer_3":"Nice"}`,
					},
				},
			},
		})
	}))
	defer server.Close()

	q, err := generateQuestionWithProvider(context.Background(), aiProviderConfig{
		Name:     "openrouter",
		APIKey:   "router-key",
		Model:    "openai/gpt-5.1-chat",
		Endpoint: server.URL,
		Headers: map[string]string{
			"HTTP-Referer": "http://localhost:3000",
			"X-Title":      "Mindmeld",
		},
		SourceTag: "openrouter",
	}, contentrating.Work, "geography")
	if err != nil {
		t.Fatalf("generateQuestionWithProvider returned error: %v", err)
	}

	if q.Source != "openrouter" {
		t.Fatalf("expected source openrouter, got %q", q.Source)
	}
	if q.QuestionText == "" || q.CorrectAnswer == "" {
		t.Fatalf("expected populated question, got %+v", q)
	}
	if seenAuthorization != "Bearer router-key" {
		t.Fatalf("expected bearer auth header, got %q", seenAuthorization)
	}
	if seenReferer != "http://localhost:3000" {
		t.Fatalf("expected referer header, got %q", seenReferer)
	}
	if seenTitle != "Mindmeld" {
		t.Fatalf("expected title header, got %q", seenTitle)
	}
	if seenModel != "openai/gpt-5.1-chat" {
		t.Fatalf("expected model to be forwarded, got %q", seenModel)
	}
	if seenResponseFormat["type"] != "json_schema" {
		t.Fatalf("expected json_schema response format, got %#v", seenResponseFormat)
	}
	if len(seenPlugins) != 1 || seenPlugins[0]["id"] != "response-healing" {
		t.Fatalf("expected response-healing plugin, got %#v", seenPlugins)
	}
	if !strings.Contains(seenSystemPrompt, "Mode 2: stated fact.") {
		t.Fatalf("expected stated-fact guidance in system prompt, got %q", seenSystemPrompt)
	}
	if !strings.Contains(seenSystemPrompt, "Never fabricate personal facts.") {
		t.Fatalf("expected anti-fabrication guidance in system prompt, got %q", seenSystemPrompt)
	}
	if seenUserPrompt != "Now generate a question from this input:\ngeography" {
		t.Fatalf("unexpected user prompt %q", seenUserPrompt)
	}
}

func TestBuildQuestionAssistPromptsForPersonalShell(t *testing.T) {
	systemPrompt, userPrompt := buildQuestionAssistPrompts(contentrating.Work, "A personal question about Justin that I can answer")

	if !strings.Contains(systemPrompt, "Mode 3: personal question shell.") {
		t.Fatalf("expected personal-shell guidance in system prompt, got %q", systemPrompt)
	}
	if !strings.Contains(systemPrompt, "set correct_answer to exactly [fill in correct answer].") {
		t.Fatalf("expected placeholder guidance in system prompt, got %q", systemPrompt)
	}
	if userPrompt != "Now generate a question from this input:\nA personal question about Justin that I can answer" {
		t.Fatalf("unexpected user prompt %q", userPrompt)
	}
}

func TestBuildQuestionAssistPromptsForFirstPersonFact(t *testing.T) {
	systemPrompt, userPrompt := buildQuestionAssistPrompts(contentrating.Work, "my favorite fruit is blueberry")

	if !strings.Contains(systemPrompt, "rewrite it using [MY_NAME] as the subject placeholder") {
		t.Fatalf("expected [MY_NAME] guidance in system prompt, got %q", systemPrompt)
	}
	if userPrompt != "Now generate a question from this input:\nFirst-person stated fact about [MY_NAME]: my favorite fruit is blueberry" {
		t.Fatalf("unexpected user prompt %q", userPrompt)
	}
}

func TestGenerateLocalQuestionFromFirstPersonFact(t *testing.T) {
	q := generateLocalQuestion(contentrating.Work, "my favorite fruit is blueberry")

	if q.Source != "local-fallback" {
		t.Fatalf("expected local fallback source, got %q", q.Source)
	}
	if q.QuestionText != "What is [MY_NAME]'s favorite fruit?" {
		t.Fatalf("unexpected question text %q", q.QuestionText)
	}
	if q.CorrectAnswer != "blueberry" {
		t.Fatalf("unexpected correct answer %q", q.CorrectAnswer)
	}
	answers := []string{q.CorrectAnswer, q.WrongAnswer1, q.WrongAnswer2, q.WrongAnswer3}
	seen := make(map[string]bool, len(answers))
	for _, answer := range answers {
		key := strings.ToLower(answer)
		if seen[key] {
			t.Fatalf("expected unique answers, got %+v", answers)
		}
		seen[key] = true
	}
}

func TestCleanTopicAllowsLongerPromptInputs(t *testing.T) {
	input := "Make a personal multiple-choice question from this fact: Jess' go-to karaoke song is Mr. Brightside and keep Jess in the question text."
	cleaned := cleanTopic(input)

	if cleaned != input {
		t.Fatalf("expected topic to remain intact, got %q", cleaned)
	}
	if len(cleaned) > maxAssistInputLen {
		t.Fatalf("expected cleaned topic to be at most %d chars, got %d", maxAssistInputLen, len(cleaned))
	}
}
