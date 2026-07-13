package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vinhala/diestimme/backend/internal/analysis"
	"github.com/vinhala/diestimme/backend/internal/domain"
)

type CRMReader interface {
	ListCompanies(ctx context.Context) ([]domain.Company, error)
	ListOpportunities(ctx context.Context, companyID string) ([]domain.Opportunity, error)
}

type AnalysisRunner interface {
	Run(ctx context.Context, recording analysis.Recording) (domain.AnalysisResultResponse, error)
}

type Handler struct {
	crm              CRMReader
	analysisRunner   AnalysisRunner
	logger           *slog.Logger
	maxUploadBytes   int64
	allowedMediaType map[string]bool
}

func NewRouter(crm CRMReader, analysisRunner AnalysisRunner, logger *slog.Logger, maxUploadBytes int64, allowedMediaType map[string]bool) *gin.Engine {
	if logger == nil {
		logger = slog.Default()
	}

	handler := Handler{
		crm:              crm,
		analysisRunner:   analysisRunner,
		logger:           logger,
		maxUploadBytes:   maxUploadBytes,
		allowedMediaType: allowedMediaType,
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/companies", handler.listCompanies)
	api.GET("/companies/:companyId/opportunities", handler.listOpportunities)
	api.POST("/analyses", handler.createAnalysis)

	return router
}

func (h Handler) listCompanies(c *gin.Context) {
	companies, err := h.crm.ListCompanies(c.Request.Context())
	if err != nil {
		h.serverError(c, "companies could not be loaded", err)
		return
	}

	c.JSON(http.StatusOK, companies)
}

func (h Handler) listOpportunities(c *gin.Context) {
	companyID := strings.TrimSpace(c.Param("companyId"))
	if companyID == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "company id is required"})
		return
	}

	opportunities, err := h.crm.ListOpportunities(c.Request.Context(), companyID)
	if err != nil {
		h.serverError(c, "opportunities could not be loaded", err)
		return
	}

	c.JSON(http.StatusOK, opportunities)
}

func (h Handler) createAnalysis(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.maxUploadBytes+1024*1024)
	if err := c.Request.ParseMultipartForm(h.maxUploadBytes + 1024*1024); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "multipart form data could not be read: " + err.Error()})
		return
	}

	companyID := strings.TrimSpace(c.PostForm("companyId"))
	opportunityID := strings.TrimSpace(c.PostForm("opportunityId"))
	if companyID == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "company id is required"})
		return
	}
	if opportunityID == "" {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "opportunity id is required"})
		return
	}

	fileHeader, err := c.FormFile("recording")
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "recording is required"})
		return
	}
	if fileHeader.Size > h.maxUploadBytes {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: fmt.Sprintf("recording exceeds %d bytes", h.maxUploadBytes)})
		return
	}

	contentType, err := h.validateRecording(fileHeader.Filename, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Error: "recording could not be opened"})
		return
	}
	defer file.Close()

	result, err := h.analysisRunner.Run(c.Request.Context(), analysis.Recording{
		CompanyID:     companyID,
		OpportunityID: opportunityID,
		File:          file,
		Filename:      fileHeader.Filename,
		ContentType:   contentType,
	})
	if err != nil {
		h.serverError(c, "analysis could not be completed", err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h Handler) validateRecording(filename, rawContentType string) (string, error) {
	extension := strings.ToLower(filepath.Ext(filename))
	if extension != ".mp3" && extension != ".mp4" {
		return "", fmt.Errorf("recording must be an MP3 or MP4 file")
	}

	contentType := rawContentType
	if parsed, _, err := mime.ParseMediaType(rawContentType); err == nil {
		contentType = parsed
	}
	if contentType == "" || contentType == "application/octet-stream" {
		switch extension {
		case ".mp3":
			contentType = "audio/mpeg"
		case ".mp4":
			contentType = "audio/mp4"
		}
	}
	if !h.allowedMediaType[contentType] {
		return "", fmt.Errorf("recording media type %q is not supported", contentType)
	}

	return contentType, nil
}

func (h Handler) serverError(c *gin.Context, message string, err error) {
	h.logger.Error(message, "error", err)
	c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Error: err.Error()})
}
