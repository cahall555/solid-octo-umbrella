package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
	"strings"

	"solid-octo-umbrella/actions"
	"solid-octo-umbrella/models"
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
	msg := models.Message{
		Role:    "user",
		Content: chat,
	}
	req := models.Request{
		Model:    "llama3.2",
		Messages: []models.Message{msg},
		Stream:   false,
	}
	resp, err := actions.ChatOllama(localOllamaChatURL, req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Model: %s\n", resp.Model)
	fmt.Printf("Response: %s\n", resp.Message.Content)
	fmt.Printf("Completed in %v", time.Since(start))
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
