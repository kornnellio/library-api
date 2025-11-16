package transport

import (
	"net/http"
	"strconv"

	"library-api/internal/errors"
	"library-api/internal/logger"
	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Handler holds all HTTP handlers and dependencies
type Handler struct {
	userService *service.UserService
	bookService *service.BookService
}

// NewHandler creates a new handler instance
func NewHandler(userService *service.UserService, bookService *service.BookService) *Handler {
	return &Handler{
		userService: userService,
		bookService: bookService,
	}
}

// Home handles the home endpoint
func (h *Handler) Home(c *gin.Context) {
	c.String(http.StatusOK, "Library API - Running on Cloud Run")
}

// Health handles the health check endpoint
func (h *Handler) Health(c *gin.Context) {
	// Health check should verify database connectivity
	// For now, just return OK - can be enhanced to check DB via service
	c.String(http.StatusOK, "OK")
}

// Register handles user registration
func (h *Handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid registration request", zap.Error(err))
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.userService.Register(ctx, &req); err != nil {
		log.Error("Registration failed", zap.Error(err), zap.String("email", req.Email))
		respondError(c, err)
		return
	}

	log.Info("User registered successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, gin.H{"message": "Registered successfully"})
}

// Login handles user login
func (h *Handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn("Invalid login request", zap.Error(err))
		respondError(c, errors.ErrBadRequest)
		return
	}

	user, err := h.userService.Login(ctx, &req)
	if err != nil {
		log.Warn("Login failed", zap.String("email", req.Email))
		respondError(c, err)
		return
	}

	log.Info("User logged in successfully", zap.Int("user_id", user.ID))
	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user_id": user.ID})
}

// GetBooks handles retrieving all books
func (h *Handler) GetBooks(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	books, err := h.bookService.GetAll(ctx)
	if err != nil {
		log.Error("Failed to retrieve books", zap.Error(err))
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, books)
}

// CreateBook handles creating a new book
func (h *Handler) CreateBook(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Warn("Invalid book creation request", zap.Error(err))
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.bookService.Create(ctx, &book); err != nil {
		log.Error("Failed to create book", zap.Error(err))
		respondError(c, err)
		return
	}

	log.Info("Book created", zap.Int("book_id", book.ID))
	c.JSON(http.StatusCreated, book)
}

// UpdateBook handles updating an existing book
func (h *Handler) UpdateBook(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	idStr := c.Param("id")
	bookID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("Invalid book ID", zap.String("id", idStr))
		respondError(c, errors.ErrBadRequest)
		return
	}

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		log.Warn("Invalid book update request", zap.Error(err))
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.bookService.Update(ctx, bookID, &book); err != nil {
		log.Error("Failed to update book", zap.Error(err), zap.Int("book_id", bookID))
		respondError(c, err)
		return
	}

	log.Info("Book updated", zap.Int("book_id", bookID))
	c.JSON(http.StatusOK, gin.H{"message": "Book updated"})
}

// DeleteBook handles deleting a book
func (h *Handler) DeleteBook(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLogger(c)

	idStr := c.Param("id")
	bookID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Warn("Invalid book ID", zap.String("id", idStr))
		respondError(c, errors.ErrBadRequest)
		return
	}

	if err := h.bookService.Delete(ctx, bookID); err != nil {
		log.Error("Failed to delete book", zap.Error(err), zap.Int("book_id", bookID))
		respondError(c, err)
		return
	}

	log.Info("Book deleted", zap.Int("book_id", bookID))
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

// respondError handles error responses consistently
func respondError(c *gin.Context, err error) {
	var appErr *errors.AppError
	if e, ok := err.(*errors.AppError); ok {
		appErr = e
	} else {
		appErr = errors.ErrInternalServerError
	}

	c.JSON(appErr.Code, gin.H{
		"error": appErr.Message,
	})
}
