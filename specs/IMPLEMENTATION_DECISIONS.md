# DieStimme Implementation Decisions Report

## Context

The implementation turns DieStimme from a scaffold into the synchronous proof-of-concept described in `EPICS.md` and `ARC.md`: a Nuxt web interface uploads a requirements session recording to a Go backend, the backend uses OpenAI to transcribe and analyze it, and the result is stored in Twenty CRM as a note linked to the selected opportunity.

## Backend Architecture

The Go backend was implemented with Gin, matching the architecture decision that the backend should be a lightweight Gin application server. The code is split into focused internal packages:

- `handlers`: HTTP routes, request validation, upload handling, and JSON error responses.
- `twenty`: Twenty CRM REST client.
- `openaiapi`: raw HTTP client for OpenAI transcription and questionnaire generation.
- `analysis`: orchestration of transcription, questionnaire generation, note creation, and note linking.
- `questionnaire`: embedded standardized questionnaire template.
- `domain` and `config`: shared API models and environment-driven configuration.

This split keeps external API concerns away from HTTP handlers and makes the core workflow testable without real OpenAI or Twenty calls.

## Synchronous PoC Flow

The analysis endpoint is synchronous: `POST /api/analyses` waits for transcription, questionnaire generation, note creation, and note-target linking before returning.

This was chosen because the requested scope is a proof of concept and the final user flow requires displaying the generated questionnaire immediately. A queued design would need persistent job state and polling endpoints, which would add infrastructure beyond the current architecture.

The response type was changed from a queued start response to a completed result:

- `analysisId`
- `status: "completed"`
- `noteId`
- `questionnaireMarkdown`

## OpenAI Integration

OpenAI calls are implemented with raw HTTP rather than the Go SDK. This keeps the integration explicit and avoids coupling the PoC to SDK abstractions for newer speech/response parameters.

Transcription uses:

- Endpoint: `/v1/audio/transcriptions`
- Default model: `gpt-4o-transcribe-diarize`
- `response_format=diarized_json`
- `chunking_strategy=auto`

The diarized segment response is converted into speaker-labelled dialogue. If segments are not returned, the implementation falls back to the response `text` field.

Questionnaire generation uses:

- Endpoint: `/v1/responses`
- Default model: `gpt-4.1-mini`
- A prompt that preserves the questionnaire structure, returns only Markdown, fills only transcript-supported answers, and marks missing information as open clarification questions.

Both model names are configurable with `OPENAI_TRANSCRIPTION_MODEL` and `OPENAI_ANALYSIS_MODEL`.

## Questionnaire Source

The standardized onboarding questionnaire was added as an embedded Markdown file at `backend/internal/questionnaire/client_onboarding.md`.

Embedding the file gives the PoC a versioned questionnaire without requiring runtime file mounts or deploy-time environment payloads. It is still easy to review and update as a normal Markdown artifact.

## Twenty CRM Integration

Twenty is accessed through its generated Core REST API:

- Companies: `GET /rest/companies`
- Opportunities: `GET /rest/opportunities` filtered by `companyId`
- Notes: `POST /rest/notes`
- Note linking: `POST /rest/noteTargets`

The implementation normalizes `TWENTY_API_URL` so both `http://host` and `http://host/rest` work. This makes local and deployment configuration less brittle.

The filled questionnaire is stored as a Twenty note using `bodyV2.markdown`, then linked to the selected opportunity through a `noteTarget` with `targetOpportunityId`. This matches Twenty’s standard note-target relationship model.

## Upload Validation

The backend enforces MP3/MP4 uploads and limits recordings to 25 MB by default. The previous 100 MB default was reduced because OpenAI’s file upload transcription path currently limits audio files to 25 MB.

The maximum is still configurable through `MAX_UPLOAD_BYTES`, and the frontend keeps lightweight extension validation for faster feedback while relying on the backend for authoritative enforcement.

## Frontend Decisions

The existing Nuxt page was kept as a simple operational form rather than redesigned. The UI now:

- Loads companies and company-scoped opportunities from the backend.
- Shows “Analyzing...” while the synchronous backend request runs.
- Displays the final Markdown questionnaire after success.
- Shows the Twenty note ID returned by the backend.
- Preserves concrete backend error messages through Nuxt server route proxy handling.

The shared TypeScript API type was updated to `AnalysisResultResponse` so the frontend contract matches the backend response.

## Error Handling

Backend errors consistently return JSON as `{ "error": "..." }`.

The Nuxt server proxy preserves backend error text in `data.error`, allowing the browser UI to show concrete OpenAI, Twenty, upload, or validation failures instead of generic proxy errors.

## Testing And Verification

Backend tests were added for:

- Twenty response parsing and error propagation.
- Opportunity filtering by company.
- Note creation and note-target linking.
- Analysis orchestration order.
- Multipart upload success and validation failures.

Verification completed:

- Go tests passed via Docker with `golang:1.24-bookworm`.
- Nuxt typecheck passed.
- Nuxt production build passed.

The local shell did not have `go` or `gofmt`, so Go formatting and tests were run through Docker instead.

## Known Limitations

- The analysis flow is synchronous and may take a while for longer recordings.
- There is no persistent job history outside the final CRM note.
- OpenAI and Twenty calls are not retried.
- The questionnaire template is a first version and may need domain refinement after real FDE feedback.
- End-to-end live testing requires valid OpenAI and Twenty credentials plus reachable CRM data.
