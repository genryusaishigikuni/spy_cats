package mission

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler for the mission domain
type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	missionGroup := r.Group("/missions")
	{
		missionGroup.POST("", h.createMission)       // POST /missions
		missionGroup.GET("", h.listMissions)         // GET /missions
		missionGroup.GET("/:id", h.getMissionByID)   // GET /missions/:id
		missionGroup.DELETE("/:id", h.deleteMission) // DELETE /missions/:id

		// Assign cat to an existing mission
		missionGroup.PATCH("/:id/assign-cat/:catId", h.assignCat)
		missionGroup.PATCH("/:id/complete", h.markMissionComplete)
	}

	// Mark a Target as complete
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

// listMissions handles GET /missions
func (h *Handler) listMissions(c *gin.Context) {
	missions, err := h.service.ListMissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, missions)
}

// getMissionByID handles GET /missions/:id
func (h *Handler) getMissionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	m, err := h.service.GetMissionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, m)
}

// deleteMission handles DELETE /missions/:id
func (h *Handler) deleteMission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	if err := h.service.DeleteMission(uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// assignCat handles PATCH /missions/:id/assign-cat/:catId
func (h *Handler) assignCat(c *gin.Context) {
	missionIDStr := c.Param("id")
	catIDStr := c.Param("catId")

	missionID, err := strconv.Atoi(missionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}
	catID, err := strconv.Atoi(catIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	if err := h.service.AssignCat(uint(missionID), uint(catID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cat assigned successfully"})
}

// markMissionComplete handles PATCH /missions/:id/complete
func (h *Handler) markMissionComplete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mission ID"})
		return
	}

	if err := h.service.MarkMissionComplete(uint(id)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mission marked as completed"})
}
