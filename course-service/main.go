package main

import (
	"log"

	"final/course-service/config"
	"final/course-service/handlers"
	"final/course-service/middleware"

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

	// Public routes
	r.GET("/courses", handlers.GetCourses)
	r.GET("/courses/:id", handlers.GetCourse)

	// Admin only routes
	admin := r.Group("/")
	admin.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/courses/:id/students", handlers.GetCourseStudents)
		admin.POST("/courses", handlers.CreateCourse)
		admin.PUT("/courses/:id", handlers.UpdateCourse)
		admin.DELETE("/courses/:id", handlers.DeleteCourse)
	}

	log.Println("Course Service running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatal(err)
	}
}
