package note

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the "note" domain.
type Handler struct {
	service Service
}

// NewHandler creates a new note Handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up note endpoints, for example to create/update notes.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/targets/:targetId/notes", h.createNote)
	r.PUT("/notes/:id", h.updateNote)
}

// createNote handles POST /targets/:targetId/notes
func (h *Handler) createNote(c *gin.Context) {
	targetIDStr := c.Param("targetId")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n, err := h.service.CreateNote(uint(targetID), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, n)
}

// updateNote handles PUT /notes/:id
func (h *Handler) updateNote(c *gin.Context) {
	idStr := c.Param("id")
	noteID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNote, err := h.service.UpdateNote(uint(noteID), req.Content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedNote)
}
