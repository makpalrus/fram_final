package handlers

import (
	"final/enrollment-service/client"
	"final/enrollment-service/config"
	"final/enrollment-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateEnrollment(c *gin.Context) {
	var req models.CreateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userData, err := client.GetUserByID(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	courseData, err := client.GetCourseByID(req.CourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course not found"})
		return
	}

	// Check duplicate enrollment
	var existing models.Enrollment
	if err := config.DB.Where("user_id = ? AND course_id = ?", req.UserID, req.CourseID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Already enrolled in this course"})
		return
	}

	enrollment := models.Enrollment{
		UserID:   req.UserID,
		CourseID: req.CourseID,
		Status:   "active",
	}

	if err := config.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create enrollment"})
		return
	}

	response := models.EnrollmentResponse{
		Enrollment: enrollment,
		User: models.UserSummary{
			ID:    uint(userData["id"].(float64)),
			Name:  userData["name"].(string),
			Email: userData["email"].(string),
		},
		Course: models.CourseSummary{
			ID:    uint(courseData["id"].(float64)),
			Title: courseData["title"].(string),
			Price: courseData["price"].(float64),
		},
	}

	c.JSON(http.StatusCreated, response)
}

func GetEnrollments(c *gin.Context) {
	var enrollments []models.Enrollment
	config.DB.Find(&enrollments)
	c.JSON(http.StatusOK, enrollments)
}

func GetEnrollment(c *gin.Context) {
	var enrollment models.Enrollment
	if err := config.DB.First(&enrollment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}
	c.JSON(http.StatusOK, enrollment)
}

func GetUserEnrollments(c *gin.Context) {
	var enrollments []models.Enrollment
	userID := c.Param("user_id")
	config.DB.Where("user_id = ?", userID).Find(&enrollments)
	c.JSON(http.StatusOK, enrollments)
}

func GetCourseEnrollments(c *gin.Context) {
	courseID := c.Param("course_id")
	var enrollments []models.Enrollment
	config.DB.Where("course_id = ?", courseID).Find(&enrollments)
	c.JSON(http.StatusOK, gin.H{
		"course_id":     courseID,
		"student_count": len(enrollments),
		"enrollments":   enrollments,
	})
}

func DeleteEnrollment(c *gin.Context) {
	var enrollment models.Enrollment
	if err := config.DB.First(&enrollment, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}
	config.DB.Delete(&enrollment)
	c.JSON(http.StatusOK, gin.H{"message": "Enrollment cancelled successfully"})
}
