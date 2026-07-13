package analysis

import (
	"context"
	"io"
	"strings"
	"testing"
)

type fakeTranscriber struct {
	called bool
}

func (f *fakeTranscriber) Transcribe(ctx context.Context, file io.Reader, filename, contentType string) (string, error) {
	f.called = true
	if filename != "meeting.mp3" || contentType != "audio/mpeg" {
		return "", errString("unexpected recording metadata")
	}
	return "Speaker 1: hello", nil
}

type fakeGenerator struct {
	calledWithTranscript string
}

func (f *fakeGenerator) GenerateQuestionnaire(ctx context.Context, questionnaire, transcript string) (string, error) {
	if !strings.Contains(questionnaire, "Template") {
		return "", errString("missing questionnaire")
	}
	f.calledWithTranscript = transcript
	return "# Filled", nil
}

type fakeProvider struct{}

func (fakeProvider) Load() string {
	return "# Template"
}

type fakeCRM struct {
	noteMarkdown string
	targetNoteID string
	targetOppID  string
}

func (f *fakeCRM) CreateNote(ctx context.Context, title, markdown string) (string, error) {
	f.noteMarkdown = markdown
	return "note-1", nil
}

func (f *fakeCRM) CreateNoteTarget(ctx context.Context, noteID, opportunityID string) error {
	f.targetNoteID = noteID
	f.targetOppID = opportunityID
	return nil
}

func TestServiceRunCompletesAnalysisAndPersistsNote(t *testing.T) {
	transcriber := &fakeTranscriber{}
	generator := &fakeGenerator{}
	crm := &fakeCRM{}
	service := NewService(transcriber, generator, fakeProvider{}, crm)

	result, err := service.Run(context.Background(), Recording{
		OpportunityID: "opp-1",
		File:          strings.NewReader("audio"),
		Filename:      "meeting.mp3",
		ContentType:   "audio/mpeg",
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if result.Status != "completed" || result.NoteID != "note-1" || result.QuestionnaireMarkdown != "# Filled" {
		t.Fatalf("unexpected result %#v", result)
	}
	if result.AnalysisID == "" {
		t.Fatalf("expected analysis id")
	}
	if !transcriber.called || generator.calledWithTranscript != "Speaker 1: hello" {
		t.Fatalf("pipeline did not call transcription and generation")
	}
	if crm.noteMarkdown != "# Filled" || crm.targetNoteID != "note-1" || crm.targetOppID != "opp-1" {
		t.Fatalf("unexpected crm calls %#v", crm)
	}
}

type errString string

func (e errString) Error() string {
	return string(e)
}
