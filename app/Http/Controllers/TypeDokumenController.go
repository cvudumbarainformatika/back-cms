package controllers

import (
	"net/http"
	"strconv"

	models "github.com/cvudumbarainformatika/backend/app/Models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type TypeDokumenController struct {
	DB *sqlx.DB
}

func NewTypeDokumenController(db *sqlx.DB) *TypeDokumenController {
	return &TypeDokumenController{
		DB: db,
	}
}

// GET LIST
func (ctl *TypeDokumenController) GetList(c *gin.Context) {

	data, total, err := models.GetAllTypeDokumen(
		ctl.DB,
		map[string]interface{}{},
		0,
		100,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  data,
		"total": total,
	})
}

// GET BY ID
func (ctl *TypeDokumenController) GetByID(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	data, err := models.FindTypeDokumenByID(ctl.DB, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

// CREATE
func (ctl *TypeDokumenController) Create(c *gin.Context) {

	var data models.TypeDokumen

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err := data.Create(ctl.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}

// UPDATE
func (ctl *TypeDokumenController) Update(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	data, err := models.FindTypeDokumenByID(ctl.DB, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	if err := c.ShouldBindJSON(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = data.Update(ctl.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    data,
	})
}

// DELETE
func (ctl *TypeDokumenController) Delete(c *gin.Context) {

	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	data, err := models.FindTypeDokumenByID(ctl.DB, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "data not found",
		})
		return
	}

	err = data.Delete(ctl.DB)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
