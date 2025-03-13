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
		Name  string `json:"name"`
		Breed string `json:"breed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Uses the service layer to create the cat
	cat, err := h.service.CreateCat(req.Name, req.Breed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Responds with the created cat
	c.JSON(http.StatusCreated, cat)
}

// listCats handles GET /cats
func (h *Handler) listCats(c *gin.Context) {
	// Uses the service layer to get all cats
	cats, err := h.service.ListCats()
	if err != nil {
		// If there's an error, returns a 500 internal server error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Returns the list of cats in the response
	c.JSON(http.StatusOK, cats)
}

// getCat handles GET /cats/:id
func (h *Handler) getCat(c *gin.Context) {
	idStr := c.Param("id") // Gets the ID from the URL parameters

	// Convert the string ID to uint
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	// Calls the service layer to get the cat by its ID
	cat, err := h.service.GetCat(uint(id))
	if err != nil {
		// If the cat is not found or another error occurs
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	// Return the found cat as a JSON response
	c.JSON(http.StatusOK, cat)
}

// updateCat handles PUT /cats/:id
func (h *Handler) updateCat(c *gin.Context) {
	idStr := c.Param("id") // Gets the ID from the URL parameters

	// Convert the string ID to uint
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	// Binds the request body to a struct to get the new name and breed
	var req struct {
		Name  string `json:"name"`
		Breed string `json:"breed"`
	}

	// Binds the JSON payload to the struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Calls the service to update the cat
	updatedCat, err := h.service.UpdateCat(uint(id), req.Name, req.Breed)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	// Returns the updated cat as a JSON response
	c.JSON(http.StatusOK, updatedCat)
}

// deleteCat handles DELETE /cats/:id
func (h *Handler) deleteCat(c *gin.Context) {
	idStr := c.Param("id") // Get the ID from the URL parameters

	// Convert the string ID to uint
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cat ID"})
		return
	}

	// Calls the service to delete the cat
	err = h.service.DeleteCat(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cat not found"})
		return
	}

	// Returns 204 No Content on successful deletion
	c.Status(http.StatusNoContent)
}
