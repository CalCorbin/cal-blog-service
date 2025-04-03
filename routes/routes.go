// routes/routes.go
package routes

import (
	"cal-blog-service/controllers"
	"cal-blog-service/middleware"
	"cal-blog-service/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// SetupRouter initializes the Gin router and defines all routes
func SetupRouter(db *gorm.DB) *gin.Engine {
	engine := gin.Default()

	if err := engine.SetTrustedProxies(nil); err != nil {
		log.Printf("Uh oh, failed to set trusted proxies: %v", err)
	}

	// Initialize controllers
	authController := &controllers.AuthController{DB: db}

	// Auth routes (public)
	engine.POST("/auth/register", authController.Register)
	engine.POST("/auth/login", authController.Login)

	// API Routes
	engine.GET("/posts", func(context *gin.Context) {
		var posts []models.BlogPost
		db.Find(&posts)
		context.JSON(http.StatusOK, posts)
	})

	engine.GET("/posts/:id", func(context *gin.Context) {
		id := context.Param("id")
		var post models.BlogPost

		if err := db.First(&post, id).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}

		context.JSON(http.StatusOK, post)
	})

	// Protected API Routes - require auth
	authorized := engine.Group("/")
	authorized.Use(middleware.AuthMiddleware())
	{
		// Create post - requires authentication
		authorized.POST("/posts", func(context *gin.Context) {
			var post models.BlogPost
			if err := context.BindJSON(&post); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Set author ID to current user
			userID, _ := context.Get("userID")
			post.AuthorID = userID.(uint)

			db.Create(&post)
			context.JSON(http.StatusCreated, post)
		})

		// Update post - requires authentication and ownership
		authorized.PUT("/posts/:id", func(context *gin.Context) {
			id := context.Param("id")
			var post models.BlogPost

			if err := db.First(&post, id).Error; err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
				return
			}

			// Check ownership or admin role
			userID, _ := context.Get("userID")
			userRole, _ := context.Get("userRole")

			if post.AuthorID != userID.(uint) && userRole != "admin" {
				context.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this post"})
				return
			}

			// Only update fields from JSON input
			var updatedPost models.BlogPost
			if err := context.BindJSON(&updatedPost); err != nil {
				context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			// Don't allow changing the author
			updatedPost.ID = post.ID
			updatedPost.AuthorID = post.AuthorID

			db.Save(&updatedPost)
			context.JSON(http.StatusOK, updatedPost)
		})

		// Delete post - requires authentication and ownership
		authorized.DELETE("/posts/:id", func(context *gin.Context) {
			id := context.Param("id")
			var post models.BlogPost

			if err := db.First(&post, id).Error; err != nil {
				context.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
				return
			}

			// Check ownership or admin role
			userID, _ := context.Get("userID")
			userRole, _ := context.Get("userRole")

			if post.AuthorID != userID.(uint) && userRole != "admin" {
				context.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this post"})
				return
			}

			db.Delete(&post)
			context.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
		})
	}

	return engine
}
