package cat

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the "cat" domain.
type Handler struct {
	service Service
}

// NewHandler creates a new cat Handler.
func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

// RegisterRoutes sets up the cat endpoints under "/cats".
func (h *Handler) RegisterRoutes(r *gin.Engine) {
	catGroup := r.Group("/cats")
	{
		catGroup.POST("", h.createCat)       // POST /cats
		catGroup.GET("", h.listCats)         // GET /cats
		catGroup.GET("/:id", h.getCat)       // GET /cats/:id
		catGroup.PUT("/:id", h.updateCat)    // PUT /cats/:id
		catGroup.DELETE("/:id", h.deleteCat) // DELETE /cats/:id
	}
}

// createCat handles POST /cats
func (h *Handler) createCat(c *gin.Context) {
	var req struct {
		Name              string  `json:"name"`
		Breed             string  `json:"breed"`
		YearsOfExperience int     `json:"years_of_experience"`
		Salary            float64 `json:"salary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.service.CreateCat(req.Name, req.Breed, req.YearsOfExperience, req.Salary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// listCats handles GET /cats
func (h *Handler) listCats(c *gin.Context) {
	cats, err := h.service.ListCats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// getCat handles GET /cats/:id
func (h *Handler) getCat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	cat, err := h.service.GetCat(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

// updateCat handles PUT /cats/:id
func (h *Handler) updateCat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	var req struct {
		Name              string  `json:"name"`
		Breed             string  `json:"breed"`
		YearsOfExperience int     `json:"years_of_experience"`
		Salary            float64 `json:"salary"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	updatedCat, err := h.service.UpdateCat(uint(id), req.Name, req.Breed, req.YearsOfExperience, req.Salary)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedCat)
}

// deleteCat handles DELETE /cats/:id
func (h *Handler) deleteCat(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	err = h.service.DeleteCat(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
