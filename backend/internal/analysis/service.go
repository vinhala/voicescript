package analysis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/vinhala/voicescript/backend/internal/domain"
)

type Transcriber interface {
	Transcribe(ctx context.Context, file io.Reader, filename, contentType string) (string, error)
}

type QuestionnaireGenerator interface {
	GenerateQuestionnaire(ctx context.Context, questionnaire, transcript string) (string, error)
}

type QuestionnaireProvider interface {
	Load() string
}

type CRM interface {
	CreateNote(ctx context.Context, title, markdown string) (string, error)
	CreateNoteTarget(ctx context.Context, noteID, opportunityID string) error
}

type Service struct {
	transcriber            Transcriber
	questionnaireGenerator QuestionnaireGenerator
	questionnaireProvider  QuestionnaireProvider
	crm                    CRM
}

type Recording struct {
	CompanyID     string
	OpportunityID string
	File          io.Reader
	Filename      string
	ContentType   string
}

func NewService(transcriber Transcriber, questionnaireGenerator QuestionnaireGenerator, questionnaireProvider QuestionnaireProvider, crm CRM) *Service {
	return &Service{
		transcriber:            transcriber,
		questionnaireGenerator: questionnaireGenerator,
		questionnaireProvider:  questionnaireProvider,
		crm:                    crm,
	}
}

func (s *Service) Run(ctx context.Context, recording Recording) (domain.AnalysisResultResponse, error) {
	transcript, err := s.transcriber.Transcribe(ctx, recording.File, recording.Filename, recording.ContentType)
	if err != nil {
		return domain.AnalysisResultResponse{}, err
	}

	questionnaire := s.questionnaireProvider.Load()
	if questionnaire == "" {
		return domain.AnalysisResultResponse{}, fmt.Errorf("questionnaire template is empty")
	}

	markdown, err := s.questionnaireGenerator.GenerateQuestionnaire(ctx, questionnaire, transcript)
	if err != nil {
		return domain.AnalysisResultResponse{}, err
	}

	noteID, err := s.crm.CreateNote(ctx, "voicescript requirements analysis", markdown)
	if err != nil {
		return domain.AnalysisResultResponse{}, err
	}
	if err := s.crm.CreateNoteTarget(ctx, noteID, recording.OpportunityID); err != nil {
		return domain.AnalysisResultResponse{}, err
	}

	return domain.AnalysisResultResponse{
		AnalysisID:            newAnalysisID(),
		Status:                "completed",
		NoteID:                noteID,
		QuestionnaireMarkdown: markdown,
	}, nil
}

func newAnalysisID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "analysis"
	}

	return hex.EncodeToString(bytes[:])
}
