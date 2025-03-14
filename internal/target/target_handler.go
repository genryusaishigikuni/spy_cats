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

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the target-related endpoints.
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	// POST /missions/:missionId/targets to add a target to a specific mission

	// DELETE /targets/:id to remove a target by its ID
	r.DELETE("/targets/:id", h.removeTarget)
}

// addTarget handles POST /missions/:missionId/targets

// removeTarget handles DELETE /targets/:id
func (h *Handler) removeTarget(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target ID"})
		return
	}

	if err := h.service.RemoveTarget(uint(id)); err != nil {
		// If the error is "cannot delete a completed target", we return a 403 Forbidden
		if err.Error() == "cannot delete a completed target" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		// Handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Target removed successfully"})
}
