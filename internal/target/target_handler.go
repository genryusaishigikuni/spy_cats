package target

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the "target" domain.
type Handler struct {
	service Service
}

// NewHandler creates a new target Handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the target-related endpoints.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// POST /missions/:missionId/targets to add a target to a specific mission
	r.POST("/missions/:missionId/targets", h.addTarget)

	// DELETE /targets/:id to remove a target by its ID
	r.DELETE("/targets/:id", h.removeTarget)
}

// addTarget handles POST /missions/:missionId/targets
func (h *Handler) addTarget(c *gin.Context) {
	missionIDStr := c.Param("missionId")
	missionID, err := strconv.Atoi(missionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddTarget(uint(missionID), req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Target added successfully"})
}

// removeTarget handles DELETE /targets/:id
func (h *Handler) removeTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.service.RemoveTarget(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target removed successfully"})
}
