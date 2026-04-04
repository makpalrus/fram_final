package models

type Enrollment struct {
	ID       uint `json:"id" gorm:"primaryKey"`
	UserID   uint `json:"user_id"`
	CourseID uint `json:"course_id"`
}
