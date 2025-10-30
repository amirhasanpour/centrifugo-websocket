package centrifugo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/centrifugal/centrifuge-go"
)

type CentrifugoClient struct {
	client *centrifuge.Client
	apiURL string
	apiKey string
}

type PublishRequest struct {
	Channel string      `json:"channel"`
	Data    any         `json:"data"`
}

type PublishResponse struct {
	Error  *centrifuge.Error 		 `json:"error,omitempty"`
	Result *centrifuge.PublishResult `json:"result,omitempty"`
}

func NewCentrifugoClient() (*CentrifugoClient, error) {
	centrifugoURL := os.Getenv("CENTRIFUGO_URL")
	if centrifugoURL == "" {
		centrifugoURL = "http://localhost:8000"
	}

	apiKey := os.Getenv("CENTRIFUGO_API_KEY")
	if apiKey == "" {
		apiKey = "your-api-key"
	}

	config := centrifuge.Config{}

	// The WebSocket endpoint is passed here
	client := centrifuge.NewJsonClient(centrifugoURL+"/connection/websocket", config)

	// Optional: set token if needed
	// client.SetToken("your-token-here")

	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to Centrifugo: %w", err)
	}

	log.Println("Successfully connected to Centrifugo")

	return &CentrifugoClient{
		client: client,
		apiURL: centrifugoURL,
		apiKey: apiKey,
	}, nil
}

func (c *CentrifugoClient) Publish(channel string, data any) error {
	// Using HTTP API for publishing
	apiURL := c.apiURL + "/api/publish"

	publishReq := PublishRequest{
		Channel: channel,
		Data:    data,
	}

	jsonData, err := json.Marshal(publishReq)
	if err != nil {
		return fmt.Errorf("failed to marshal publish request: %w", err)
	}

	form := url.Values{}
	form.Add("data", string(jsonData))

	req, err := http.NewRequest("POST", apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = form.Encode()
	req.Header.Set("Authorization", "apikey "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	httpClient := &http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send publish request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("publish failed with status: %d", resp.StatusCode)
	}

	log.Printf("Published message to channel: %s", channel)
	return nil
}

func (c *CentrifugoClient) Close() {
	c.client.Close()
}

func (c *CentrifugoClient) GenerateToken(userID string, exp int64) (string, error) {
	// In a real implementation, you would generate a JWT token
	// For now, we'll use a simple approach or rely on Centrifugo's built-in auth
	// This would typically involve creating a JWT with the Centrifugo claims
	return "", nil
}