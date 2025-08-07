package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"murim-helper/internal/domain"
)

type OllamaService interface {
	GenerateScheduleFromText(description string) ([]domain.Schedule, error)
}

type ollamaService struct{}

func NewOllamaService() OllamaService {
	return &ollamaService{}
}

func (s *ollamaService) GenerateScheduleFromText(description string) ([]domain.Schedule, error) {
	prompt := fmt.Sprintf(`
You are a discipline assistant. Based on this input: "%s",
generate a full-day schedule in structured JSON format. 

User wakes up at 07:00 and sleeps at 23:00. Include Bible reading in the morning and prayer at night.
Add meals, work, rest, gym, and any mentioned custom items.

Format:
[
  {
    "title": "Task Title",
    "description": "What it is",
    "start_time": "HH:MM",
    "end_time": "HH:MM"
  }
]
Return ONLY valid JSON array.
`, description)

	reqData := map[string]string{
		"model":  "phi3",
		"prompt": prompt,
	}
	jsonData, _ := json.Marshal(reqData)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var rawResp struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &rawResp); err != nil {
		return nil, fmt.Errorf("ollama response decode failed: %w", err)
	}

	// 🛠 Extract JSON array between [ and ]
	raw := rawResp.Response
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start == -1 || end == -1 || end <= start {
		return nil, fmt.Errorf("no JSON array found in response")
	}

	jsonOnly := raw[start : end+1]

	return domain.ParseSchedulesFromJSON(jsonOnly)
}
