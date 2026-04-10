package handlers

import (
	"final/config"
	"final/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEnrollment(c *gin.Context) {
	var enrollment models.Enrollment
	if err := c.ShouldBindJSON(&enrollment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Create(&enrollment)
	c.JSON(http.StatusCreated, enrollment)
}
func GetEnrollments(c *gin.Context) {
	var enrollments []models.Enrollment
	if err := config.DB.Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных"})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}
func GetUserEnrollments(c *gin.Context) {
	var enrollments []models.Enrollment
	userID := c.Param("user_id")
	config.DB.Where("user_id = ?", userID).Find(&enrollments)
	c.JSON(http.StatusOK, enrollments)
}

func UpdateEnrollment(c *gin.Context) {
	var enrollment models.Enrollment
	if err := config.DB.First(&enrollment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}
	c.ShouldBindJSON(&enrollment)
	config.DB.Save(&enrollment)
	c.JSON(http.StatusOK, enrollment)
}

func DeleteEnrollment(c *gin.Context) {
	var enrollment models.Enrollment
	if err := config.DB.First(&enrollment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}
	config.DB.Delete(&enrollment)
	c.JSON(http.StatusOK, gin.H{"message": "Enrollment deleted"})
}
