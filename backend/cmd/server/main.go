package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/vinhala/voicescript/backend/internal/analysis"
	"github.com/vinhala/voicescript/backend/internal/config"
	"github.com/vinhala/voicescript/backend/internal/handlers"
	"github.com/vinhala/voicescript/backend/internal/openaiapi"
	"github.com/vinhala/voicescript/backend/internal/questionnaire"
	"github.com/vinhala/voicescript/backend/internal/twenty"
)

func main() {
	cfg := config.FromEnv()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	httpClient := &http.Client{Timeout: 10 * time.Minute}
	crmClient := twenty.NewClient(cfg.TwentyAPIURL, cfg.TwentyAPIKey, httpClient)
	openAIClient := openaiapi.NewClient(
		cfg.OpenAIAPIKey,
		cfg.OpenAITranscriptionModel,
		cfg.OpenAIAnalysisModel,
		httpClient,
	)
	questionnaireProvider := questionnaire.NewProvider()
	analysisService := analysis.NewService(openAIClient, openAIClient, questionnaireProvider, crmClient)
	router := handlers.NewRouter(crmClient, analysisService, logger, cfg.MaxUploadBytes, cfg.AllowedMediaType)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger.Info("starting voicescript backend", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("backend stopped", "error", err)
		os.Exit(1)
	}
}
