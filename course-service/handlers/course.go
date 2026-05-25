package handlers

import (
	"final/course-service/config"
	"final/course-service/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func GetCourses(c *gin.Context) {
	var courses []models.Course
	search := c.Query("search")
	if search != "" {
		config.DB.Where("title ILIKE ?", fmt.Sprintf("%%%s%%", search)).Find(&courses)
	} else {
		config.DB.Find(&courses)
	}
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

func GetCourseStudents(c *gin.Context) {
	id := c.Param("id")
	enrollmentServiceURL := "http://enrollment-service:8083"

	var result interface{}
	resp, err := resty.New().R().
		SetResult(&result).
		Get(fmt.Sprintf("%s/enrollments/course/%s", enrollmentServiceURL, id))

	if err != nil || resp.IsError() {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Could not fetch students"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create course"})
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
	if err := config.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}
	c.JSON(http.StatusOK, course)
}

func DeleteCourse(c *gin.Context) {
	var course models.Course
	if err := config.DB.First(&course, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}
