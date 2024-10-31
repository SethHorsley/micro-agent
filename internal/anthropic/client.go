package anthropic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	baseURL    = "https://api.anthropic.com/v1/messages"
	apiVersion = "2023-06-01"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	Stream    bool      `json:"stream"`
}

type ContentBlock struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type StreamResponse struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
	Content []ContentBlock `json:"content"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(apiKey string) *Client {
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}
	return &Client{
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

func parseSSEMessage(data string) (*StreamResponse, error) {
	if !strings.HasPrefix(data, "{") {
		return nil, nil // Not JSON data
	}

	var resp StreamResponse
	if err := json.Unmarshal([]byte(data), &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *Client) SendMessage(model string, message string, onChunk func(string)) error {
	req := Request{
		Model:     model,
		MaxTokens: 1024,
		Messages:  []Message{{Role: "user", Content: message}},
		Stream:    true,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	request, err := http.NewRequest("POST", baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	request.Header.Set("x-api-key", c.apiKey)
	request.Header.Set("anthropic-version", apiVersion)
	request.Header.Set("content-type", "application/json")

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response (status %d): %s", resp.StatusCode, string(body))
	}

	scanner := bufio.NewScanner(resp.Body)
	var messageBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines
		if line == "" {
			continue
		}

		// Handle SSE format
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Handle stream end
			if data == "[DONE]" {
				break
			}

			// Parse the JSON data
			streamResp, err := parseSSEMessage(data)
			if err != nil {
				continue // Skip malformed messages
			}

			// Skip non-JSON messages
			if streamResp == nil {
				continue
			}

			// Handle different message types
			switch streamResp.Type {
			case "content_block_delta":
				if streamResp.Delta.Text != "" {
					messageBuilder.WriteString(streamResp.Delta.Text)
					onChunk(streamResp.Delta.Text)
				}
			case "message_stop":
				return nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}
