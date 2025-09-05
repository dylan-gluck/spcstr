package handlers

import (
	"log/slog"
	"strings"
)

func HandleUserPrompt(projectRoot string, data map[string]interface{}) error {
	prompt, _ := data["prompt"].(string)
	slog.Debug("user_prompt_submit", "prompt_length", len(prompt))

	if containsSensitiveData(prompt) {
		slog.Warn("prompt may contain sensitive data")
	}

	return nil
}

func containsSensitiveData(text string) bool {
	sensitive := []string{
		"password",
		"secret",
		"token",
		"api_key",
		"apikey",
		"private_key",
		"privatekey",
		"ssh_key",
		"sshkey",
	}

	lower := strings.ToLower(text)
	for _, s := range sensitive {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
