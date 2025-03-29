// routes/routes.go
package routes

import (
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

	engine.POST("/posts", func(context *gin.Context) {
		var post models.BlogPost
		if err := context.BindJSON(&post); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&post)
		context.JSON(http.StatusCreated, post)
	})

	engine.PUT("/posts/:id", func(context *gin.Context) {
		id := context.Param("id")
		var post models.BlogPost
		if err := db.First(&post, id).Error; err != nil {
			context.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if err := context.BindJSON(&post); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&post)
		context.JSON(http.StatusOK, post)
	})

	engine.DELETE("/posts/:id", func(context *gin.Context) {
		id := context.Param("id")
		db.Delete(&models.BlogPost{}, id)
		context.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	})

	return engine
}
