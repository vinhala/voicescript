package domain

type AnalysisResultResponse struct {
	AnalysisID            string `json:"analysisId"`
	Status                string `json:"status"`
	NoteID                string `json:"noteId"`
	QuestionnaireMarkdown string `json:"questionnaireMarkdown"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
