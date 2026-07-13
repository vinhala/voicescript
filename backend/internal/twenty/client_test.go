package twenty

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCompaniesParsesTwentyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/companies" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Fatalf("missing auth header")
		}
		if r.URL.Query().Get("depth") != "0" || r.URL.Query().Get("limit") != "100" {
			t.Fatalf("unexpected query %s", r.URL.RawQuery)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"companies": []map[string]string{
					{"id": "company-1", "name": "Acme"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL+"/rest", "test-key", server.Client())
	companies, err := client.ListCompanies(context.Background())
	if err != nil {
		t.Fatalf("ListCompanies returned error: %v", err)
	}
	if len(companies) != 1 || companies[0].ID != "company-1" || companies[0].Name != "Acme" {
		t.Fatalf("unexpected companies %#v", companies)
	}
}

func TestListOpportunitiesFiltersByCompany(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/rest/opportunities" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("filter"); got != `companyId[in]:["company-1"]` {
			t.Fatalf("unexpected filter %q", got)
		}

		_ = json.NewEncoder(w).Encode(map[string]any{
			"data": map[string]any{
				"opportunities": []map[string]string{
					{"id": "opp-1", "name": "Rollout", "companyId": "company-1"},
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", server.Client())
	opportunities, err := client.ListOpportunities(context.Background(), "company-1")
	if err != nil {
		t.Fatalf("ListOpportunities returned error: %v", err)
	}
	if len(opportunities) != 1 || opportunities[0].ID != "opp-1" || opportunities[0].CompanyID != "company-1" {
		t.Fatalf("unexpected opportunities %#v", opportunities)
	}
}

func TestCreateNoteAndTarget(t *testing.T) {
	var sawNote bool
	var sawTarget bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest/notes":
			sawNote = true
			var body struct {
				Title  string `json:"title"`
				BodyV2 struct {
					Markdown string `json:"markdown"`
				} `json:"bodyV2"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode note body: %v", err)
			}
			if body.Title != "Title" || body.BodyV2.Markdown != "# Markdown" {
				t.Fatalf("unexpected note body %#v", body)
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"createNote": map[string]string{"id": "note-1"}},
			})
		case "/rest/noteTargets":
			sawTarget = true
			var body map[string]string
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decode note target body: %v", err)
			}
			if body["noteId"] != "note-1" || body["targetOpportunityId"] != "opp-1" {
				t.Fatalf("unexpected note target body %#v", body)
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"data": map[string]any{"createNoteTarget": map[string]string{"id": "target-1"}},
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", server.Client())
	noteID, err := client.CreateNote(context.Background(), "Title", "# Markdown")
	if err != nil {
		t.Fatalf("CreateNote returned error: %v", err)
	}
	if noteID != "note-1" {
		t.Fatalf("unexpected note id %q", noteID)
	}
	if err := client.CreateNoteTarget(context.Background(), noteID, "opp-1"); err != nil {
		t.Fatalf("CreateNoteTarget returned error: %v", err)
	}
	if !sawNote || !sawTarget {
		t.Fatalf("expected both note and note target requests")
	}
}

func TestTwentyErrorIncludesMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error":    "BadRequestException",
			"messages": []string{"invalid filter"},
		})
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key", server.Client())
	_, err := client.ListCompanies(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	if got := err.Error(); got != "twenty request failed: invalid filter" {
		t.Fatalf("unexpected error %q", got)
	}
}
