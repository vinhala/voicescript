package twenty

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/vinhala/voicescript/backend/internal/domain"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	normalizedBaseURL := strings.TrimRight(baseURL, "/")
	normalizedBaseURL = strings.TrimSuffix(normalizedBaseURL, "/rest")

	return &Client{
		baseURL:    strings.TrimRight(normalizedBaseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c *Client) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	endpoint, err := c.restURL("/companies", url.Values{
		"depth":    {"0"},
		"limit":    {"100"},
		"order_by": {"name[AscNullsLast]"},
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			Companies []domain.Company `json:"companies"`
		} `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &response); err != nil {
		return nil, err
	}

	return response.Data.Companies, nil
}

func (c *Client) ListOpportunities(ctx context.Context, companyID string) ([]domain.Opportunity, error) {
	endpoint, err := c.restURL("/opportunities", url.Values{
		"depth":    {"0"},
		"limit":    {"100"},
		"filter":   {fmt.Sprintf("companyId[in]:[%q]", companyID)},
		"order_by": {"name[AscNullsLast]"},
	})
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			Opportunities []domain.Opportunity `json:"opportunities"`
		} `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodGet, endpoint, nil, &response); err != nil {
		return nil, err
	}

	return response.Data.Opportunities, nil
}

func (c *Client) CreateNote(ctx context.Context, title, markdown string) (string, error) {
	endpoint, err := c.restURL("/notes", url.Values{"depth": {"0"}})
	if err != nil {
		return "", err
	}

	requestBody := map[string]any{
		"title": title,
		"bodyV2": map[string]any{
			"markdown":  markdown,
			"blocknote": nil,
		},
	}
	var response struct {
		Data struct {
			CreateNote struct {
				ID string `json:"id"`
			} `json:"createNote"`
		} `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodPost, endpoint, requestBody, &response); err != nil {
		return "", err
	}
	if response.Data.CreateNote.ID == "" {
		return "", fmt.Errorf("twenty returned a note without an id")
	}

	return response.Data.CreateNote.ID, nil
}

func (c *Client) CreateNoteTarget(ctx context.Context, noteID, opportunityID string) error {
	endpoint, err := c.restURL("/noteTargets", url.Values{"depth": {"0"}})
	if err != nil {
		return err
	}

	requestBody := map[string]any{
		"noteId":              noteID,
		"targetOpportunityId": opportunityID,
	}
	var response struct {
		Data struct {
			CreateNoteTarget struct {
				ID string `json:"id"`
			} `json:"createNoteTarget"`
		} `json:"data"`
	}
	if err := c.doJSON(ctx, http.MethodPost, endpoint, requestBody, &response); err != nil {
		return err
	}
	if response.Data.CreateNoteTarget.ID == "" {
		return fmt.Errorf("twenty returned a note target without an id")
	}

	return nil
}

func (c *Client) restURL(path string, values url.Values) (string, error) {
	if c.baseURL == "" {
		return "", fmt.Errorf("TWENTY_API_URL is required")
	}

	parsed, err := url.Parse(c.baseURL + "/rest" + path)
	if err != nil {
		return "", err
	}
	parsed.RawQuery = values.Encode()

	return parsed.String(), nil
}

func (c *Client) doJSON(ctx context.Context, method, endpoint string, body any, target any) error {
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(payload)
	}

	request, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	request.Header.Set("Accept", "application/json")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		request.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("twenty request failed: %s", errorMessage(responseBody, response.Status))
	}
	if err := json.Unmarshal(responseBody, target); err != nil {
		return fmt.Errorf("decode twenty response: %w", err)
	}

	return nil
}

func errorMessage(body []byte, fallback string) string {
	var payload struct {
		Error    string   `json:"error"`
		Message  string   `json:"message"`
		Messages []string `json:"messages"`
	}
	if err := json.Unmarshal(body, &payload); err == nil {
		switch {
		case len(payload.Messages) > 0:
			return strings.Join(payload.Messages, "; ")
		case payload.Message != "":
			return payload.Message
		case payload.Error != "":
			return payload.Error
		}
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed != "" {
		return trimmed
	}
	return fallback
}
