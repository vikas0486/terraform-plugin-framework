package provider

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ErrNotFound is returned when the API responds with a 404,
// signaling that the resource no longer exists upstream.
var ErrNotFound = errors.New("keystore not found")

// APIClient holds configuration required to communicate
// with the Thales REST API.
type APIClient struct {
	Endpoint string
	Client   *http.Client
}

// NewClient creates a reusable API client.
func NewClient(endpoint string) *APIClient {
	return &APIClient{
		Endpoint: endpoint,
		Client:   &http.Client{},
	}
}

// CreateKeystore creates a new keystore via REST API.
func (c *APIClient) CreateKeystore(name string) (*Keystore, error) {

	body := map[string]string{
		"name": name,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		c.Endpoint+"/keystores",
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d creating keystore: %s", resp.StatusCode, string(b))
	}

	var keystore Keystore

	if err := json.NewDecoder(resp.Body).Decode(&keystore); err != nil {
		return nil, err
	}

	return &keystore, nil
}

// GetKeystore retrieves an existing keystore.
func (c *APIClient) GetKeystore(id string) (*Keystore, error) {

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s/keystores/%s", c.Endpoint, id),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d reading keystore: %s", resp.StatusCode, string(b))
	}

	var keystore Keystore

	if err := json.NewDecoder(resp.Body).Decode(&keystore); err != nil {
		return nil, err
	}

	return &keystore, nil
}
