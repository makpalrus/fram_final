package main

import (
	"final/config"
	"final/handlers"
	"final/middleware"
	"final/models"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	config.DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Enrollment{})

	r := gin.Default()

	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", handlers.GetUsers)
		protected.POST("/users", handlers.CreateUser)
		protected.GET("/users/:id", handlers.GetUser)
		protected.PUT("/users/:id", handlers.UpdateUser)
		protected.DELETE("/users/:id", handlers.DeleteUser)

		protected.GET("/courses", handlers.GetCourses)
		protected.POST("/courses", handlers.CreateCourse)
		protected.GET("/courses/:id", handlers.GetCourse)
		protected.PUT("/courses/:id", handlers.UpdateCourse)
		protected.DELETE("/courses/:id", handlers.DeleteCourse)

		protected.POST("/enrollments", handlers.CreateEnrollment)
		protected.GET("/enrollments/:user_id", handlers.GetUserEnrollments)
		protected.PUT("/enrollments/:id", handlers.UpdateEnrollment)
		protected.DELETE("/enrollments/:id", handlers.DeleteEnrollment)
		protected.GET("/enrollments", handlers.GetEnrollments)
	}

	r.Run(":8080")
}
