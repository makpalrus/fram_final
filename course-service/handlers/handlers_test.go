package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"final/course-service/config"
	"final/course-service/models"

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
	if err := db.AutoMigrate(&models.Course{}); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}
	config.DB = db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/courses", GetCourses)
	r.GET("/courses/:id", GetCourse)
	r.POST("/courses", CreateCourse)
	r.PUT("/courses/:id", UpdateCourse)
	r.DELETE("/courses/:id", DeleteCourse)
	return r
}

func seedCourse(t *testing.T, title, description string, price float64) models.Course {
	t.Helper()
	c := models.Course{Title: title, Description: description, Price: price, Duration: 10}
	if err := config.DB.Create(&c).Error; err != nil {
		t.Fatalf("failed to seed course: %v", err)
	}
	return c
}

// TEST 1: Create course returns 201
func TestCreateCourse_Success(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	body := `{"title":"Go Basics","description":"Intro to Go","price":49.99,"duration":10}`
	req, _ := http.NewRequest("POST", "/courses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp models.Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Go Basics", resp.Title)
}

// TEST 2: Create course with missing title returns 400
func TestCreateCourse_BadJSON(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	body := `{"description":"no title","price":10}`
	req, _ := http.NewRequest("POST", "/courses", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TEST 3: Get all courses returns list
func TestGetCourses_ReturnsList(t *testing.T) {
	setupTestDB(t)
	seedCourse(t, "Python Basics", "Learn Python", 29.99)
	seedCourse(t, "Docker Deep Dive", "Containers", 39.99)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/courses", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var courses []models.Course
	json.Unmarshal(w.Body.Bytes(), &courses)
	assert.Equal(t, 2, len(courses))
}

// TEST 4: Search courses by name (LIKE)
func TestGetCourses_Search(t *testing.T) {
	setupTestDB(t)
	seedCourse(t, "Golang Advanced", "Advanced Go", 59.99)
	seedCourse(t, "Python Basics", "Learn Python", 19.99)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/courses?search=golang", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var courses []models.Course
	json.Unmarshal(w.Body.Bytes(), &courses)
	assert.Equal(t, 1, len(courses))
	assert.Equal(t, "Golang Advanced", courses[0].Title)
}

// TEST 5: Get course by ID returns course
func TestGetCourse_Found(t *testing.T) {
	setupTestDB(t)
	c := seedCourse(t, "Kubernetes", "K8s fundamentals", 79.99)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/courses/"+itoa(c.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Kubernetes", resp.Title)
}

// TEST 6: Get course by invalid ID returns 404
func TestGetCourse_NotFound(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/courses/9999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TEST 7: Update course returns 200
func TestUpdateCourse_Success(t *testing.T) {
	setupTestDB(t)
	c := seedCourse(t, "Old Title", "Old desc", 10.0)
	r := setupRouter()

	body := `{"title":"New Title","description":"New desc","price":20.0,"duration":5}`
	req, _ := http.NewRequest("PUT", "/courses/"+itoa(c.ID), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp models.Course
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "New Title", resp.Title)
}

// TEST 8: Update non-existent course returns 404
func TestUpdateCourse_NotFound(t *testing.T) {
	setupTestDB(t)
	r := setupRouter()

	body := `{"title":"X","description":"X","price":1.0,"duration":1}`
	req, _ := http.NewRequest("PUT", "/courses/9999", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TEST 9: Delete course returns 200
func TestDeleteCourse_Success(t *testing.T) {
	setupTestDB(t)
	c := seedCourse(t, "To Delete", "desc", 5.0)
	r := setupRouter()

	req, _ := http.NewRequest("DELETE", "/courses/"+itoa(c.ID), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TEST 10: Search with no match returns empty list
func TestGetCourses_SearchNoMatch(t *testing.T) {
	setupTestDB(t)
	seedCourse(t, "React Basics", "Learn React", 24.99)
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/courses?search=blockchain", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var courses []models.Course
	json.Unmarshal(w.Body.Bytes(), &courses)
	assert.Equal(t, 0, len(courses))
}

func itoa(id uint) string {
	return fmt.Sprintf("%d", id)
}
