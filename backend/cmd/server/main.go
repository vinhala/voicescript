package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/vinhala/diestimme/backend/internal/analysis"
	"github.com/vinhala/diestimme/backend/internal/config"
	"github.com/vinhala/diestimme/backend/internal/handlers"
	"github.com/vinhala/diestimme/backend/internal/openaiapi"
	"github.com/vinhala/diestimme/backend/internal/questionnaire"
	"github.com/vinhala/diestimme/backend/internal/twenty"
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

	logger.Info("starting DieStimme backend", "port", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("backend stopped", "error", err)
		os.Exit(1)
	}
}
