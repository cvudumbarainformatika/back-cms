package controllers

import (
	"net/http"
	"strconv"

	"github.com/cvudumbarainformatika/backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// ExampleController handles example requests
// This is a template controller showing best practices
type ExampleController struct {
	db *sqlx.DB
}

// NewExampleController creates a new ExampleController instance
func NewExampleController(db *sqlx.DB) *ExampleController {
	return &ExampleController{
		db: db,
	}
}

// GetAll handles GET /api/v1/examples with pagination
func (ec *ExampleController) GetAll(c *gin.Context) {
	// Get pagination parameters from query string
	params := utils.GetFilterParams(c)

	// TODO: Implement your GetAll logic
	// Example:
	// items, total, err := models.GetExamples(ec.db, params)
	// if err != nil {
	//     utils.Error(c, http.StatusInternalServerError, "internal_error", "Failed to get examples", nil)
	//     return
	// }
	//
	// pagination := utils.CreateLaravelPagination(c, items, params.Page, params.PerPage, total)
	// c.JSON(http.StatusOK, pagination)

	utils.Error(c, http.StatusNotImplemented, "not_implemented", "This endpoint is not implemented yet", nil)
}

// GetByID handles GET /api/v1/examples/:id
func (ec *ExampleController) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid example ID format", nil)
		return
	}

	// TODO: Implement your GetByID logic
	// Example:
	// item, err := models.FindExampleByID(ec.db, id)
	// if err != nil {
	//     utils.Error(c, http.StatusInternalServerError, "internal_error", "Failed to get example", nil)
	//     return
	// }
	//
	// if item == nil {
	//     utils.Error(c, http.StatusNotFound, "example_not_found", "Example not found", nil)
	//     return
	// }
	//
	// utils.Success(c, http.StatusOK, "Example retrieved successfully", gin.H{
	//     "example": item,
	// })

	utils.Error(c, http.StatusNotImplemented, "not_implemented", "This endpoint is not implemented yet", nil)
}

// Create handles POST /api/v1/examples
func (ec *ExampleController) Create(c *gin.Context) {
	// TODO: Implement your Create logic
	// var req requests.CreateExampleRequest
	//
	// if err := req.Validate(c); err != nil {
	//     return  // Validation errors are handled by Validate()
	// }
	//
	// item := &models.Example{
	//     // Map request fields to model
	// }
	//
	// if err := item.Create(ec.db); err != nil {
	//     utils.Error(c, http.StatusInternalServerError, "internal_error", "Failed to create example", nil)
	//     return
	// }
	//
	// utils.Success(c, http.StatusCreated, "Example created successfully", gin.H{
	//     "example": item,
	// })

	utils.Error(c, http.StatusNotImplemented, "not_implemented", "This endpoint is not implemented yet", nil)
}

// Update handles PUT /api/v1/examples/:id
func (ec *ExampleController) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid example ID format", nil)
		return
	}

	// TODO: Implement your Update logic
	// 1. Fetch existing item
	// 2. Validate request
	// 3. Update fields
	// 4. Save to database
	// 5. Return updated item

	utils.Error(c, http.StatusNotImplemented, "not_implemented", "This endpoint is not implemented yet", nil)
}

// Delete handles DELETE /api/v1/examples/:id
func (ec *ExampleController) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "invalid_id", "Invalid example ID format", nil)
		return
	}

	// TODO: Implement your Delete logic
	// 1. Check if item exists
	// 2. Delete from database
	// 3. Return success response

	utils.Error(c, http.StatusNotImplemented, "not_implemented", "This endpoint is not implemented yet", nil)
}
