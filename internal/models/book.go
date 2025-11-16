package models

// Book represents a book in the library
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}
