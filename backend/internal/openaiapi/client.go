package openaiapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type Client struct {
	apiKey             string
	baseURL            string
	transcriptionModel string
	analysisModel      string
	httpClient         *http.Client
}

func NewClient(apiKey, transcriptionModel, analysisModel string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		apiKey:             apiKey,
		baseURL:            "https://api.openai.com/v1",
		transcriptionModel: transcriptionModel,
		analysisModel:      analysisModel,
		httpClient:         httpClient,
	}
}

func (c *Client) Transcribe(ctx context.Context, file io.Reader, filename, contentType string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is required")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(filename)))
	if contentType != "" {
		partHeader.Set("Content-Type", contentType)
	}
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}
	if err := writer.WriteField("model", c.transcriptionModel); err != nil {
		return "", err
	}
	if err := writer.WriteField("response_format", "diarized_json"); err != nil {
		return "", err
	}
	if err := writer.WriteField("chunking_strategy", "auto"); err != nil {
		return "", err
	}
	if err := writer.Close(); err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/audio/transcriptions", &body)
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Accept", "application/json")

	var response transcriptionResponse
	if err := c.do(request, &response); err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	transcript := response.Dialogue()
	if transcript == "" {
		return "", fmt.Errorf("transcription returned no text")
	}

	return transcript, nil
}

func (c *Client) GenerateQuestionnaire(ctx context.Context, questionnaire, transcript string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is required")
	}

	payload := map[string]any{
		"model": c.analysisModel,
		"input": []map[string]string{
			{
				"role": "system",
				"content": strings.Join([]string{
					"You fill out Voiceline client onboarding questionnaires from requirements elicitation transcripts.",
					"Return only markdown.",
					"Preserve the questionnaire headings and bullet structure.",
					"Fill answers only when supported by the transcript.",
					"Where information is missing, write a concise open clarification question.",
				}, " "),
			},
			{
				"role":    "user",
				"content": fmt.Sprintf("Questionnaire template:\n\n%s\n\nSpeaker-labelled transcript:\n\n%s", questionnaire, transcript),
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/responses", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+c.apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	var response responsesResponse
	if err := c.do(request, &response); err != nil {
		return "", fmt.Errorf("questionnaire analysis failed: %w", err)
	}

	markdown := strings.TrimSpace(response.Text())
	if markdown == "" {
		return "", fmt.Errorf("questionnaire analysis returned no text")
	}

	return markdown, nil
}

func (c *Client) do(request *http.Request, target any) error {
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
		return fmt.Errorf("%s", errorMessage(responseBody, response.Status))
	}
	if err := json.Unmarshal(responseBody, target); err != nil {
		return fmt.Errorf("decode openai response: %w", err)
	}

	return nil
}

type transcriptionResponse struct {
	Text     string `json:"text"`
	Segments []struct {
		Speaker string `json:"speaker"`
		Text    string `json:"text"`
	} `json:"segments"`
}

func (r transcriptionResponse) Dialogue() string {
	if len(r.Segments) == 0 {
		return strings.TrimSpace(r.Text)
	}

	lines := make([]string, 0, len(r.Segments))
	for index, segment := range r.Segments {
		text := strings.TrimSpace(segment.Text)
		if text == "" {
			continue
		}
		speaker := strings.TrimSpace(segment.Speaker)
		if speaker == "" {
			speaker = fmt.Sprintf("Speaker %d", index+1)
		}
		lines = append(lines, fmt.Sprintf("%s: %s", speaker, text))
	}

	return strings.Join(lines, "\n")
}

type responsesResponse struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}

func (r responsesResponse) Text() string {
	if strings.TrimSpace(r.OutputText) != "" {
		return r.OutputText
	}

	var builder strings.Builder
	for _, output := range r.Output {
		for _, content := range output.Content {
			if strings.TrimSpace(content.Text) == "" {
				continue
			}
			if builder.Len() > 0 {
				builder.WriteString("\n")
			}
			builder.WriteString(content.Text)
		}
	}

	return builder.String()
}

func errorMessage(body []byte, fallback string) string {
	var payload struct {
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &payload); err == nil && payload.Error.Message != "" {
		return payload.Error.Message
	}

	trimmed := strings.TrimSpace(string(body))
	if trimmed != "" {
		return trimmed
	}
	return fallback
}

func escapeQuotes(value string) string {
	return strings.ReplaceAll(value, `"`, `\"`)
}
