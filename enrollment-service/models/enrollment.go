package models

import "time"

type Enrollment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	CourseID  uint      `json:"course_id"`
	Status    string    `json:"status" gorm:"default:active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateEnrollmentRequest struct {
	UserID   uint `json:"user_id" binding:"required"`
	CourseID uint `json:"course_id" binding:"required"`
}

type EnrollmentResponse struct {
	Enrollment Enrollment    `json:"enrollment"`
	User       UserSummary   `json:"user"`
	Course     CourseSummary `json:"course"`
}

type UserSummary struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CourseSummary struct {
	ID    uint    `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}
