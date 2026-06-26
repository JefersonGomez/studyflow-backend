package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type GroqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GroqRequest struct {
	Model    string        `json:"model"`
	Messages []GroqMessage `json:"messages"`
}

type GroqChoice struct {
	Message GroqMessage `json:"message"`
}

type GroqResponse struct {
	Choices []GroqChoice `json:"choices"`
}

func Generate(prompt string) (string, error) {
	reqBody := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: []GroqMessage{
			{Role: "user", Content: prompt},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Log temporal para ver el error
	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Printf("Groq status: %d, body: %s", resp.StatusCode, string(bodyBytes))

	var groqResp GroqResponse
	if err := json.Unmarshal(bodyBytes, &groqResp); err != nil {
		return "", err
	}

	if len(groqResp.Choices) == 0 {
		return "", fmt.Errorf("no se recibió respuesta de Groq")
	}

	return groqResp.Choices[0].Message.Content, nil
}
