package ai

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func Generate(prompt string) (string, error) {

	newResques := &OllamaRequest{
		Prompt: prompt,
		Model:  "qwen2.5:7b",
		Stream: false,
	}

	body, err := json.Marshal(newResques)

	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return "", err
	}

	return ollamaResp.Response, nil

}
