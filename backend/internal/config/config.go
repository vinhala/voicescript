package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port                     string
	OpenAIAPIKey             string
	OpenAITranscriptionModel string
	OpenAIAnalysisModel      string
	TwentyAPIURL             string
	TwentyAPIKey             string
	MaxUploadBytes           int64
	AllowedMediaType         map[string]bool
}

func FromEnv() Config {
	return Config{
		Port:                     getEnv("PORT", "80"),
		OpenAIAPIKey:             os.Getenv("OPENAI_API_KEY"),
		OpenAITranscriptionModel: getEnv("OPENAI_TRANSCRIPTION_MODEL", "gpt-4o-transcribe-diarize"),
		OpenAIAnalysisModel:      getEnv("OPENAI_ANALYSIS_MODEL", "gpt-4.1-mini"),
		TwentyAPIURL:             normalizeTwentyURL(getEnv("TWENTY_API_URL", "http://twenty-server:3000")),
		TwentyAPIKey:             os.Getenv("TWENTY_API_KEY"),
		MaxUploadBytes:           getEnvInt64("MAX_UPLOAD_BYTES", 25*1024*1024),
		AllowedMediaType: map[string]bool{
			"audio/mpeg": true,
			"audio/mp3":  true,
			"audio/mp4":  true,
			"video/mp4":  true,
		},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvInt64(key string, fallback int64) int64 {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

func normalizeTwentyURL(rawURL string) string {
	normalized := strings.TrimRight(rawURL, "/")
	return strings.TrimSuffix(normalized, "/rest")
}
