package mission

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the "mission" domain.
type Handler struct {
	service Service
}

// NewHandler creates a new mission Handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the mission endpoints under "/missions", and
// any related target-completion routes, etc.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// Missions
	missionGroup := r.Group("/missions")
	{
		missionGroup.POST("", h.createMission) // POST /missions
		// I might add more mission routes as needed (e.g., GET /missions/:id, etc.)
	}

	// Marks a Target as complete
	r.PATCH("/targets/:id/complete", h.completeTarget)
}

// createMission handles POST /missions
func (h *Handler) createMission(c *gin.Context) {
	var req struct {
		CatID       uint     `json:"cat_id"`
		TargetNames []string `json:"target_names"` // minimal example
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m, err := h.service.CreateMission(req.CatID, req.TargetNames)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, m)
}

// completeTarget handles PATCH /targets/:id/complete
func (h *Handler) completeTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.service.CompleteTarget(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Target completed"})
}
