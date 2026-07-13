package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/vinhala/voicescript/backend/internal/analysis"
	"github.com/vinhala/voicescript/backend/internal/domain"
)

type fakeCRMReader struct{}

func (fakeCRMReader) ListCompanies(ctx context.Context) ([]domain.Company, error) {
	return []domain.Company{{ID: "company-1", Name: "Acme"}}, nil
}

func (fakeCRMReader) ListOpportunities(ctx context.Context, companyID string) ([]domain.Opportunity, error) {
	return []domain.Opportunity{{ID: "opp-1", Name: "Rollout", CompanyID: companyID}}, nil
}

type fakeRunner struct {
	recording analysis.Recording
}

func (f *fakeRunner) Run(ctx context.Context, recording analysis.Recording) (domain.AnalysisResultResponse, error) {
	f.recording = recording
	return domain.AnalysisResultResponse{
		AnalysisID:            "analysis-1",
		Status:                "completed",
		NoteID:                "note-1",
		QuestionnaireMarkdown: "# Filled",
	}, nil
}

func TestListCompanies(t *testing.T) {
	router := testRouter(&fakeRunner{}, 1024)
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/companies", nil)

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d", response.Code)
	}
	var companies []domain.Company
	if err := json.Unmarshal(response.Body.Bytes(), &companies); err != nil {
		t.Fatalf("decode companies: %v", err)
	}
	if len(companies) != 1 || companies[0].ID != "company-1" {
		t.Fatalf("unexpected companies %#v", companies)
	}
}

func TestCreateAnalysisSuccess(t *testing.T) {
	runner := &fakeRunner{}
	router := testRouter(runner, 1024)
	body, contentType := multipartBody(t, "company-1", "opp-1", "meeting.mp3", "audio/mpeg", []byte("audio"))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/analyses", body)
	request.Header.Set("Content-Type", contentType)

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status %d body %s", response.Code, response.Body.String())
	}
	var result domain.AnalysisResultResponse
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	if result.Status != "completed" || result.NoteID != "note-1" {
		t.Fatalf("unexpected result %#v", result)
	}
	if runner.recording.CompanyID != "company-1" || runner.recording.OpportunityID != "opp-1" {
		t.Fatalf("unexpected recording %#v", runner.recording)
	}
}

func TestCreateAnalysisRejectsMissingFields(t *testing.T) {
	router := testRouter(&fakeRunner{}, 1024)
	body, contentType := multipartBody(t, "", "opp-1", "meeting.mp3", "audio/mpeg", []byte("audio"))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/analyses", body)
	request.Header.Set("Content-Type", contentType)

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status %d", response.Code)
	}
}

func TestCreateAnalysisRejectsUnsupportedFile(t *testing.T) {
	router := testRouter(&fakeRunner{}, 1024)
	body, contentType := multipartBody(t, "company-1", "opp-1", "meeting.wav", "audio/wav", []byte("audio"))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/analyses", body)
	request.Header.Set("Content-Type", contentType)

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status %d", response.Code)
	}
}

func TestCreateAnalysisRejectsOversizedFile(t *testing.T) {
	router := testRouter(&fakeRunner{}, 5)
	body, contentType := multipartBody(t, "company-1", "opp-1", "meeting.mp3", "audio/mpeg", []byte("too large"))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/analyses", body)
	request.Header.Set("Content-Type", contentType)

	router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status %d", response.Code)
	}
}

func testRouter(runner *fakeRunner, maxUploadBytes int64) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return NewRouter(fakeCRMReader{}, runner, nil, maxUploadBytes, map[string]bool{
		"audio/mpeg": true,
		"audio/mp4":  true,
		"video/mp4":  true,
	})
}

func multipartBody(t *testing.T, companyID, opportunityID, filename, contentType string, file []byte) (*bytes.Buffer, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	if companyID != "" {
		if err := writer.WriteField("companyId", companyID); err != nil {
			t.Fatalf("write company: %v", err)
		}
	}
	if opportunityID != "" {
		if err := writer.WriteField("opportunityId", opportunityID); err != nil {
			t.Fatalf("write opportunity: %v", err)
		}
	}

	partHeader := make(textproto.MIMEHeader)
	partHeader.Set("Content-Disposition", `form-data; name="recording"; filename="`+filename+`"`)
	partHeader.Set("Content-Type", contentType)
	part, err := writer.CreatePart(partHeader)
	if err != nil {
		t.Fatalf("create file part: %v", err)
	}
	if _, err := part.Write(file); err != nil {
		t.Fatalf("write file: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}

	return body, writer.FormDataContentType()
}
