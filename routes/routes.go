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
	r := gin.Default()

	if err := r.SetTrustedProxies(nil); err != nil {
		log.Printf("Uh oh, failed to set trusted proxies: %v", err)
	}

	// API Routes
	r.GET("/posts", func(c *gin.Context) {
		var posts []models.BlogPost
		db.Find(&posts)
		c.JSON(http.StatusOK, posts)
	})

	r.POST("/posts", func(c *gin.Context) {
		var post models.BlogPost
		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Create(&post)
		c.JSON(http.StatusCreated, post)
	})

	r.PUT("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")
		var post models.BlogPost
		if err := db.First(&post, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
			return
		}
		if err := c.BindJSON(&post); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		db.Save(&post)
		c.JSON(http.StatusOK, post)
	})

	r.DELETE("/posts/:id", func(c *gin.Context) {
		id := c.Param("id")
		db.Delete(&models.BlogPost{}, id)
		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
	})

	return r
}
