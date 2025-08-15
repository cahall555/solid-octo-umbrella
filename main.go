package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"solid-octo-umbrella/actions"
	"solid-octo-umbrella/models"
	"solid-octo-umbrella/tools"
)

const localOllamaChatURL = "http://localhost:11434/api/chat"
const localOllamaGenerateURL = "http://localhost:11434/api/generate"

func main() {
	start := time.Now()

	var action string
	fmt.Println("Enter 'chat' to chat with Ollama or 'generate' to generate text:")
	fmt.Scanln(&action)
	switch action {
	case "chat":
		fmt.Println("What would you like to talk about?")
		reader := bufio.NewReader(os.Stdin)
		chat, _ := reader.ReadString('\n')
		chat = strings.TrimSpace(chat)

		chatOllama(start, chat)
	case "generate":
		fmt.Println("What would you like to generate?")
		reader := bufio.NewReader(os.Stdin)
		generate, _ := reader.ReadString('\n')
		generate = strings.TrimSpace(generate)
		generateOllama(start, generate)
	default:
		fmt.Println("Invalid action. Please enter 'chat' or 'generate'.")
		os.Exit(1)

	}
}

func chatOllama(start time.Time, chat string) {
	// Unmarshal tool JSON into struct
	var tool models.Tool
	if err := json.Unmarshal([]byte(tools.ToolGetActiveAlerts), &tool); err != nil {
		fmt.Println("tool unmarshal:", err)
		os.Exit(1)
	}

	msgs := []models.Message{
		{Role: "user", Content: chat},
	}

	req := models.Request{
		Model:    "llama3.2",
		Messages: msgs,
		Tools:    []models.Tool{tool},
		Stream:   false,
	}

	// 1st call: see if model wants the tool
	resp, err := actions.ChatOllama(localOllamaChatURL, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// If it asked to call a tool, execute it and send the result back
	if len(resp.Message.ToolCalls) > 0 {
		// Keep the assistant's message (with tool_calls) in the history
		msgs = append(msgs, resp.Message)

		for _, tc := range resp.Message.ToolCalls {
			switch tc.Function.Name {
			case "get_active_alerts":
				region, _ := tc.Function.Arguments["region"].(string)
				if region == "" {
					region = "AT"
				}
				out, err := actions.GetNOAAActiveAlerts(region)
				if err != nil {
					// Send the error back as tool content so the model can react
					msgs = append(msgs, models.Message{
						Role:    "tool",
						Name:    tc.Function.Name,
						Content: fmt.Sprintf(`{"error":%q}`, err.Error()),
					})
					continue
				}
				// Send tool result as JSON string
				b, _ := json.Marshal(out)
				msgs = append(msgs, models.Message{
					Role:    "tool",
					Name:    tc.Function.Name,
					Content: string(b),
				})
			}
		}

		// 2nd call: let the model compose a final answer using the tool output
		req.Messages = msgs
		resp, err = actions.ChatOllama(localOllamaChatURL, req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Response:\n%s\n", resp.Message.Content)
	fmt.Printf("Completed in %v\n", time.Since(start))
}

func generateOllama(start time.Time, generate string) {
	req := models.GenerateRequest{
		Model:  "llama3.2",
		Prompt: generate,
		Stream: false,
	}
	resp, err := actions.GenerateOllama(localOllamaGenerateURL, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Response: %s\n", resp.Response)
	fmt.Printf("Completed in %v", time.Since(start))
}
