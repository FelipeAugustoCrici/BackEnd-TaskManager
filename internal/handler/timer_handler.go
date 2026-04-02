package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"controltasks/internal/service"
)

type TimerHandler struct {
	svc *service.TimerService
}

func NewTimerHandler(svc *service.TimerService) *TimerHandler {
	return &TimerHandler{svc: svc}
}

// GET /api/v1/timer
func (h *TimerHandler) Get(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	timer, err := h.svc.Get(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": timer}) // null quando não existe
}

// POST /api/v1/timer
func (h *TimerHandler) Start(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	type StartTimerRequest struct {
		InitialSeconds *int `json:"initial_seconds,omitempty"`
	}

	var req StartTimerRequest
	// body é opcional — ignorar erro se vazio
	_ = c.ShouldBindJSON(&req)

	if req.InitialSeconds != nil && *req.InitialSeconds < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "initial_seconds deve ser >= 0"})
		return
	}

	initialSeconds := 0
	if req.InitialSeconds != nil {
		initialSeconds = *req.InitialSeconds
	}

	timer, err := h.svc.Start(userID, initialSeconds)
	if err != nil {
		if errors.Is(err, service.ErrTimerAlreadyActive) {
			c.JSON(http.StatusConflict, gin.H{"error": "timer já ativo"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": timer})
}

// PATCH /api/v1/timer/pause
func (h *TimerHandler) Pause(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	timer, err := h.svc.Pause(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if timer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nenhum timer ativo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": timer})
}

// PATCH /api/v1/timer/resume
func (h *TimerHandler) Resume(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	timer, err := h.svc.Resume(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if timer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nenhum timer ativo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": timer})
}

// DELETE /api/v1/timer
func (h *TimerHandler) Delete(c *gin.Context) {
	userID, ok := userIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuário não autenticado"})
		return
	}

	if err := h.svc.Delete(userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNoContent) // idempotente: já não existe
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
