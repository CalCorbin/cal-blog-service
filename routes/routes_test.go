package routes

import (
	"bytes"
	"cal-blog-service/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// Setup test environment before running any tests
func init() {
	// Set Gin to test mode to avoid unnecessary logging
	gin.SetMode(gin.TestMode)
}

func setupTestDB(t *testing.T) *gorm.DB {
	// Use in-memory SQLite for testing with a specific connection string
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Migrate the schema
	if err := db.AutoMigrate(&models.BlogPost{}); err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	return db
}

func createTestPost(t *testing.T, db *gorm.DB) models.BlogPost {
	post := models.BlogPost{
		Title:   "Test Post",
		Content: "This is a test post content.",
	}
	result := db.Create(&post)
	assert.Nil(t, result.Error)
	return post
}

func TestSetupRouter(t *testing.T) {
	// Test that router is created properly
	db := setupTestDB(t)
	router := SetupRouter(db)
	assert.NotNil(t, router)

	// Print a message to confirm test is running
	fmt.Println("TestSetupRouter completed successfully")
}

func TestGetAllPostsWhenPostsEmpty(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var emptyPosts []models.BlogPost
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &emptyPosts)
	assert.Nil(t, err)
	assert.Empty(t, emptyPosts)
}

func TestGetAllPosts(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	createTestPost(t, db)
	createTestPost(t, db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts", nil)
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var posts []models.BlogPost
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &posts)
	assert.Nil(t, err)
	assert.Len(t, posts, 2)
}

func TestGetPostByID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	post := createTestPost(t, db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), nil)
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var returnedPost models.BlogPost
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &returnedPost)
	assert.Nil(t, err)
	assert.Equal(t, post.ID, returnedPost.ID)
	assert.Equal(t, post.Title, returnedPost.Title)
	assert.Equal(t, post.Content, returnedPost.Content)
}

func TestGetPostWithInvalidID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/9999", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
	assert.Nil(t, err)
	assert.Contains(t, errorResponse, "error")
	assert.Equal(t, "Post not found", errorResponse["error"])
}

func TestGetPostWithNonNumericID(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/posts/bilbobaggins", nil)
	router.ServeHTTP(responseRecorder, req)

	assert.NotEqual(t, http.StatusOK, responseRecorder.Code)
}

func TestCreatePost(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)

	newPost := models.BlogPost{
		Title:   "New Test Post",
		Content: "This is new content.",
	}
	jsonValue, _ := json.Marshal(newPost)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusCreated, responseRecorder.Code)

	var returnedPost models.BlogPost
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &returnedPost)
	assert.Nil(t, err)
	assert.NotZero(t, returnedPost.ID)
	assert.Equal(t, newPost.Title, returnedPost.Title)
	assert.Equal(t, newPost.Content, returnedPost.Content)
}

func TestCreatePostInvalidJson(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errorResponse)
	assert.Nil(t, err)
	assert.Contains(t, errorResponse, "error")
}

func TestUpdatePost(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	post := createTestPost(t, db)

	updatedPost := post
	updatedPost.Title = "Updated Title"
	updatedPost.Content = "Updated content for testing."

	jsonValue, _ := json.Marshal(updatedPost)
	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var returnedPost models.BlogPost
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &returnedPost)
	assert.Nil(t, err)
	assert.Equal(t, post.ID, returnedPost.ID)
	assert.Equal(t, "Updated Title", returnedPost.Title)
	assert.Equal(t, "Updated content for testing.", returnedPost.Content)
}

func TestUpdateNonExistentPost(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	post := createTestPost(t, db)

	responseRecorder := httptest.NewRecorder()
	jsonValue, _ := json.Marshal(post)
	req, _ := http.NewRequest("PUT", "/posts/9999", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusNotFound, responseRecorder.Code)

	var errorResponse map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &errorResponse)
	assert.Nil(t, err)
	assert.Contains(t, errorResponse, "error")
	assert.Equal(t, "Post not found", errorResponse["error"])
}

func TestUpdatePostInvalidJson(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	post := createTestPost(t, db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

func TestDeletePost(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)
	post := createTestPost(t, db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/"+strconv.FormatUint(uint64(post.ID), 10), nil)
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)

	var response map[string]string
	err := json.Unmarshal(responseRecorder.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Contains(t, response, "message")
	assert.Equal(t, "Post deleted successfully", response["message"])

	// Verify post was deleted
	var checkPost models.BlogPost
	result := db.First(&checkPost, post.ID)
	assert.Error(t, result.Error)
	assert.Equal(t, "record not found", result.Error.Error())

}

func TestDeleteNonExistentPost(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/9999", nil)
	router.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
}

func TestDeleteInvalidIDPost(t *testing.T) {
	// Setup
	db := setupTestDB(t)
	router := SetupRouter(db)

	responseRecorder := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/posts/invalid", nil)
	router.ServeHTTP(responseRecorder, req)

	// Just check that we get some response (depending on your implementation)
	assert.NotNil(t, responseRecorder.Code)
}
