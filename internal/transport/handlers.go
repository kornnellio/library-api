package transport

import (
	"net/http"
	"strconv"
	"strings"

	"library-api/internal/models"
	"library-api/internal/repository"

	"github.com/gin-gonic/gin"
)

// Handler holds all HTTP handlers and dependencies
type Handler struct {
	userRepo *repository.UserRepository
	bookRepo *repository.BookRepository
}

// NewHandler creates a new handler instance
func NewHandler(userRepo *repository.UserRepository, bookRepo *repository.BookRepository) *Handler {
	return &Handler{
		userRepo: userRepo,
		bookRepo: bookRepo,
	}
}

// Home handles the home endpoint
func (h *Handler) Home(c *gin.Context) {
	c.String(http.StatusOK, "Library API - Running on Cloud Run ..Please do not touch :))")
}

// Health handles the health check endpoint
func (h *Handler) Health(c *gin.Context) {
	// Health check logic can be added here (e.g., DB ping)
	c.String(http.StatusOK, "OK")
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userRepo.Create(req.Email, req.Password); err != nil {
		// Check if it's a unique constraint violation (PostgreSQL error code 23505)
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "duplicate key") || strings.Contains(errStr, "23505") {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userRepo.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := h.userRepo.VerifyPassword(user.Password, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID})
}

// GetBooks handles retrieving all books
func (h *Handler) GetBooks(c *gin.Context) {
	books, err := h.bookRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// CreateBook handles creating a new book
func (h *Handler) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.bookRepo.Create(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// UpdateBook handles updating an existing book
func (h *Handler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updateErr := h.bookRepo.Update(bookID, &book); updateErr != nil {
		if updateErr.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": updateErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated"})
}

// DeleteBook handles deleting a book
func (h *Handler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	bookID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	if deleteErr := h.bookRepo.Delete(bookID); deleteErr != nil {
		if deleteErr.Error() == "book not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": deleteErr.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteErr.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}
