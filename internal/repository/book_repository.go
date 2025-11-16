package repository

import (
	"database/sql"
	"errors"

	"library-api/internal/models"
)

// BookRepository handles book data operations
type BookRepository struct {
	db *sql.DB
}

// NewBookRepository creates a new book repository
func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db: db}
}

// GetAll retrieves all books
func (r *BookRepository) GetAll() ([]models.Book, error) {
	rows, err := r.db.Query("SELECT id, title, author FROM books ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author); err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

// GetByID retrieves a book by ID
func (r *BookRepository) GetByID(id int) (*models.Book, error) {
	var book models.Book
	err := r.db.QueryRow("SELECT id, title, author FROM books WHERE id=$1", id).
		Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}
	return &book, nil
}

// Create creates a new book
func (r *BookRepository) Create(book *models.Book) error {
	err := r.db.QueryRow(
		"INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id",
		book.Title, book.Author,
	).Scan(&book.ID)
	return err
}

// Update updates an existing book
func (r *BookRepository) Update(id int, book *models.Book) error {
	result, err := r.db.Exec(
		"UPDATE books SET title=$1, author=$2 WHERE id=$3",
		book.Title, book.Author, id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}

// Delete deletes a book by ID
func (r *BookRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("book not found")
	}

	return nil
}
