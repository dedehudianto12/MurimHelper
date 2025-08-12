package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"murim-helper/internal/domain"
)

type groqService struct {
	ApiKey string
}

type GroqService interface {
	GenerateScheduleFromText(description string) ([]domain.Schedule, error)
}

func NewGroqService() GroqService {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		panic("GROQ_API_KEY environment variable is not set")
	}
	return &groqService{ApiKey: apiKey}
}

func (g *groqService) GenerateScheduleFromText(description string) ([]domain.Schedule, error) {
	prompt := fmt.Sprintf(`
	You are a discipline assistant. Based on this input: "%s",
	generate a full-day schedule in JSON format with title, description, start_time, end_time (in ISO 8601 format like "2025-08-05T07:00:00+07:00").
	Include morning Bible reading, meals, work, gym, and night prayer — even if the user doesn’t mention them.

	IMPORTANT:
	- Determine the date based on the user's description (e.g., "tomorrow", "next Monday", or specific date).
	- ALL tasks, including the default ones, must use the same inferred date.
	- If no date is mentioned, use today (%s).
	- Start the day at 07:00 and end at 23:00 in Jakarta timezone (UTC+7).

	Respond ONLY with a valid JSON array like this:
	[
		{
			"title": "Task Title",
			"description": "What it is",
			"start_time": "2025-08-05T07:00:00+07:00",
			"end_time": "2025-08-05T08:00:00+07:00"
		}
	]
	`, description, time.Now().Format("2006-01-02"))

	reqBody := map[string]interface{}{
		"model": "llama3-70b-8192",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+g.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("groq request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil || len(result.Choices) == 0 {
		return nil, fmt.Errorf("groq decode error: %w", err)
	}

	// Clean & extract only the JSON array
	content := result.Choices[0].Message.Content
	start := strings.Index(content, "[")
	end := strings.LastIndex(content, "]")
	if start == -1 || end == -1 {
		return nil, fmt.Errorf("no JSON array found in Groq response")
	}

	jsonOnly := content[start : end+1]
	return domain.ParseSchedulesFromJSON(jsonOnly)
}
