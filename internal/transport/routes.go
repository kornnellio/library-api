package transport

import "github.com/gin-gonic/gin"

// SetupRoutes configures all API routes
func (h *Handler) SetupRoutes(r *gin.Engine) {
	// Public routes
	r.GET("/", h.Home)
	r.GET("/health", h.Health)
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// Book CRUD routes (no auth yet)
	r.GET("/books", h.GetBooks)
	r.POST("/books", h.CreateBook)
	r.PUT("/books/:id", h.UpdateBook)
	r.DELETE("/books/:id", h.DeleteBook)
}
