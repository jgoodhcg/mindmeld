package trivia

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seenAuthorization = r.Header.Get("Authorization")
		seenReferer = r.Header.Get("HTTP-Referer")
		seenTitle = r.Header.Get("X-Title")

		var req chatCompletionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		seenModel = req.Model

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
}
