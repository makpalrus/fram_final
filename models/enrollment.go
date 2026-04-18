package models

import "time"

type Enrollment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	CourseID  uint      `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
}
