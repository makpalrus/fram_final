package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"final/user-service/config"
	"final/user-service/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}
	config.DB = db
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUser)
	r.PUT("/users/:id", UpdateUser)
	r.DELETE("/users/:id", DeleteUser)
	return r
}

func seedUser(t *testing.T, name, email, password string) models.User {
	t.Helper()
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := models.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
		Role:     "user",
	}
	if err := config.DB.Create(&u).Error; err != nil {
		t.Fatalf("failed to seed user: %v", err)
	}
	return u
}

// TEST 1: Register returns 201 with valid input
func TestRegister_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body := `{"name":"Alice","email":"alice@test.com","password":"secret123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Contains(t, resp["message"], "registered")
}

// TEST 2: Register returns 400 on duplicate email
func TestRegister_DuplicateEmail(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	seedUser(t, "Bob", "bob@test.com", "pass")

	body := `{"name":"Bob2","email":"bob@test.com","password":"pass2"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TEST 3: Register returns 400 on bad JSON
func TestRegister_BadJSON(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{bad json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TEST 4: Login returns 200 and JWT token
func TestLogin_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	seedUser(t, "Carol", "carol@test.com", "mypassword")

	body := `{"email":"carol@test.com","password":"mypassword"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["token"])
}

// TEST 5: Login returns 401 on wrong password
func TestLogin_WrongPassword(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	seedUser(t, "Dave", "dave@test.com", "correctpass")

	body := `{"email":"dave@test.com","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TEST 6: Login returns 401 when user doesn't exist
func TestLogin_UserNotFound(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	body := `{"email":"nobody@test.com","password":"pass"}`
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TEST 7: GetUser returns 200 with existing user
func TestGetUser_Found(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	u := seedUser(t, "Eve", "eve@test.com", "pass")

	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got models.User
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, u.Email, got.Email)
}

// TEST 8: GetUser returns 404 for missing user
func TestGetUser_NotFound(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TEST 9: UpdateUser returns 200 and persists changes
func TestUpdateUser_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	seedUser(t, "Frank", "frank@test.com", "pass")

	body := `{"name":"Frank Updated","email":"frank@test.com"}`
	req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var got models.User
	json.Unmarshal(w.Body.Bytes(), &got)
	assert.Equal(t, "Frank Updated", got.Name)
}

// TEST 10: DeleteUser returns 200 and user is gone
func TestDeleteUser_Success(t *testing.T) {
	setupTestDB(t)
	router := setupRouter()

	seedUser(t, "Grace", "grace@test.com", "pass")

	req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	req2 := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}
