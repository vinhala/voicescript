export interface Company {
  id: string
  name: string
}

export interface Opportunity {
  id: string
  name: string
  companyId: string
}

export interface AnalysisResultResponse {
  analysisId: string
  status: 'completed'
  noteId: string
  questionnaireMarkdown: string
}

export interface ErrorResponse {
  error: string
}
