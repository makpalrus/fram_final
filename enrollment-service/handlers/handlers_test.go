package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"final/enrollment-service/config"
	"final/enrollment-service/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.Enrollment{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	config.DB = db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/enrollments", GetEnrollments)
	r.GET("/enrollments/:id", GetEnrollment)
	r.GET("/enrollments/user/:user_id", GetUserEnrollments)
	r.GET("/enrollments/course/:course_id", GetCourseEnrollments)
	r.DELETE("/enrollments/:id", DeleteEnrollment)
	return r
}

func seedEnrollment(t *testing.T, userID, courseID uint) models.Enrollment {
	t.Helper()
	e := models.Enrollment{UserID: userID, CourseID: courseID, Status: "active"}
	if err := config.DB.Create(&e).Error; err != nil {
		t.Fatalf("failed to seed enrollment: %v", err)
	}
	return e
}

// TEST 1: Get all enrollments returns list
func TestGetEnrollments_ReturnsList(t *testing.T) {
	setupTestDB(t)
	seedEnrollment(t, 1, 1)
	seedEnrollment(t, 2, 1)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/enrollments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var enrollments []models.Enrollment
	json.Unmarshal(w.Body.Bytes(), &enrollments)
	assert.Equal(t, 2, len(enrollments))
}

// TEST 2: Get all enrollments when empty returns empty list
func TestGetEnrollments_Empty(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/enrollments", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var enrollments []models.Enrollment
	json.Unmarshal(w.Body.Bytes(), &enrollments)
	assert.Equal(t, 0, len(enrollments))
}

// TEST 3: Get enrollment by ID returns enrollment
func TestGetEnrollment_Found(t *testing.T) {
	setupTestDB(t)
	e := seedEnrollment(t, 1, 5)
	r := setupRouter()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/enrollments/%d", e.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Enrollment
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, uint(1), resp.UserID)
	assert.Equal(t, uint(5), resp.CourseID)
}

// TEST 4: Get enrollment by non-existent ID returns 404
func TestGetEnrollment_NotFound(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/enrollments/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TEST 5: Get enrollments by user_id returns user's enrollments
func TestGetUserEnrollments_Found(t *testing.T) {
	setupTestDB(t)
	seedEnrollment(t, 3, 1)
	seedEnrollment(t, 3, 2)
	seedEnrollment(t, 7, 1)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/enrollments/user/3", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var enrollments []models.Enrollment
	json.Unmarshal(w.Body.Bytes(), &enrollments)
	assert.Equal(t, 2, len(enrollments))
}

// TEST 6: Get enrollments by user_id with no enrollments returns empty list
func TestGetUserEnrollments_Empty(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/enrollments/user/99", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var enrollments []models.Enrollment
	json.Unmarshal(w.Body.Bytes(), &enrollments)
	assert.Equal(t, 0, len(enrollments))
}

// TEST 7: Get enrollments by course_id returns correct data
func TestGetCourseEnrollments_Found(t *testing.T) {
	setupTestDB(t)
	seedEnrollment(t, 1, 10)
	seedEnrollment(t, 2, 10)
	seedEnrollment(t, 3, 10)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/enrollments/course/10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(3), resp["student_count"])
}

// TEST 8: Delete enrollment returns 200
func TestDeleteEnrollment_Success(t *testing.T) {
	setupTestDB(t)
	e := seedEnrollment(t, 1, 1)
	r := setupRouter()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/enrollments/%d", e.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TEST 9: Delete non-existent enrollment returns 404
func TestDeleteEnrollment_NotFound(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	req, _ := http.NewRequest("DELETE", "/enrollments/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TEST 10: Deleted enrollment is no longer found
func TestDeleteEnrollment_VerifyDeleted(t *testing.T) {
	setupTestDB(t)
	e := seedEnrollment(t, 5, 5)
	r := setupRouter()

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/enrollments/%d", e.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req2, _ := http.NewRequest("GET", fmt.Sprintf("/enrollments/%d", e.ID), nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}
