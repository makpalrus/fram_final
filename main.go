package main

import (
	"final/config"
	"final/handlers"
	"final/models"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	// AutoMigrate — создаёт таблицы автоматически
	config.DB.AutoMigrate(&models.User{}, &models.Course{}, &models.Enrollment{})

	r := gin.Default()

	// Users
	r.GET("/users", handlers.GetUsers)
	r.POST("/users", handlers.CreateUser)
	r.GET("/users/:id", handlers.GetUser)
	r.PUT("/users/:id", handlers.UpdateUser)
	r.DELETE("/users/:id", handlers.DeleteUser)

	// Courses
	r.GET("/courses", handlers.GetCourses)
	r.POST("/courses", handlers.CreateCourse)
	r.GET("/courses/:id", handlers.GetCourse)
	r.PUT("/courses/:id", handlers.UpdateCourse)
	r.DELETE("/courses/:id", handlers.DeleteCourse)

	// Enrollments
	r.POST("/enrollments", handlers.CreateEnrollment)
	r.GET("/enrollments/:user_id", handlers.GetUserEnrollments)
	r.PUT("/enrollments/:id", handlers.UpdateEnrollment)
	r.DELETE("/enrollments/:id", handlers.DeleteEnrollment)

	r.Run(":8080")
}
