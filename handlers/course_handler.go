package handlers

import (
	"final/config"
	"final/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCourses(c *gin.Context) {
	var courses []models.Course
	config.DB.Find(&courses)
	c.JSON(http.StatusOK, courses)
}

func GetCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func UpdateCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	config.DB.Save(&course)
	c.JSON(http.StatusOK, course)
}

func DeleteCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	config.DB.Delete(&course)
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted"})
}
