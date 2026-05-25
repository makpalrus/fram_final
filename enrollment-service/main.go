package main

import (
	"log"

	"final/enrollment-service/config"
	"final/enrollment-service/handlers"
	"final/enrollment-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public: student can enroll and view own enrollments
	r.POST("/enrollments", handlers.CreateEnrollment)
	r.GET("/enrollments/user/:user_id", handlers.GetUserEnrollments)
	r.GET("/enrollments/course/:course_id", handlers.GetCourseEnrollments)
	r.GET("/enrollments/:id", handlers.GetEnrollment)
	r.DELETE("/enrollments/:id", handlers.DeleteEnrollment)

	// Admin only
	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/enrollments", handlers.GetEnrollments)
	}

	log.Println("Enrollment Service running on :8083")
	if err := r.Run(":8083"); err != nil {
		log.Fatal(err)
	}
}
