package service

import (
	"context"
	"fmt"
	"os"

	"murim-helper/internal/domain"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIService interface {
	GenerateScheduleFromText(desc string) ([]domain.Schedule, error)
}

type openAIService struct {
	client *openai.Client
}

func NewOpenAIService() OpenAIService {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is not set")
	}
	config := openai.DefaultConfig(apiKey)
	return &openAIService{client: openai.NewClientWithConfig(config)}
}

func (o *openAIService) GenerateScheduleFromText(desc string) ([]domain.Schedule, error) {
	prompt := fmt.Sprintf(`
	You are a discipline assistant. Based on this input: "%s",
	generate a list of structured schedule items in JSON format.

	Also the user is a software engineer who wakes up at 7:00 AM and goes to bed at 11:00 PM.
	And the user has a daily routine to read the bible in the morning and pray before going to bed.
	Generate a daily schedule for the user, including tasks like work, exercise, meals, and
	relaxation time. Each task should have a title, description, start time, and end time.

	Each item must include: title, description, start_time, end_time.
	All time should be in 24-hour "HH:MM" format.
	Output should be an array of objects like:
	[
	{
		"title": "Wake up",
		"description": "Start the day with prayer",
		"start_time": "07:00",
		"end_time": "07:30"
	}
	]
	`, desc)

	resp, err := o.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o, // use GPT-4o or GPT-3.5 if needed
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
			Temperature: 0.4,
		},
	)
	if err != nil {
		return nil, err
	}

	// Extract and parse JSON from AI response
	text := resp.Choices[0].Message.Content
	schedules, err := domain.ParseSchedulesFromJSON(text)
	if err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w", err)
	}

	return schedules, nil
}
