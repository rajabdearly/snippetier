package repo

import (
	"database/sql"
	"log"
)

type Snippet struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	UserId      int    `json:"userId"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type SnippetsRepo struct {
	db *sql.DB
}

func NewSnippetsRepo(db *sql.DB) *SnippetsRepo {
	return &SnippetsRepo{db}
}

func (r *SnippetsRepo) GetAllSnippets() ([]Snippet, error) {
	// Define the SQL query to select all rows from the "snippets" table
	query := "SELECT * FROM snippets"

	// Execute the query
	rows, err := r.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// Create a slice to store the retrieved entities
	var snippets []Snippet

	// Iterate over the result set and scan each row into a Snippet struct
	for rows.Next() {
		var snippet Snippet
		if err := rows.Scan(&snippet.ID, &snippet.Name, &snippet.Description, &snippet.Content, &snippet.UserId, &snippet.CreatedAt, &snippet.UpdatedAt); err != nil {
			log.Fatal(err)
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return snippets, nil
}

// SaveSnippet saves a single snippet to the "snippets" table using prepared statements.
func (r *SnippetsRepo) SaveSnippet(userId int, name, description, content string) (Snippet, error) {
	query := `
        INSERT INTO snippets (name, description, content, user_id)
        VALUES (?, ?, ?, ?)
    `
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return Snippet{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, content, userId)
	if err != nil {
		log.Println("Error saving snippet:", err)
		return Snippet{}, err
	}

	// Retrieve the newly created snippet'r ID from the database
	var id int
	err = r.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return Snippet{}, err
	}

	// Return the newly created snippet with the generated ID
	return Snippet{ID: id, Name: name, Description: description, Content: content}, nil
}

// UpdateSnippet updates an existing snippet in the "snippets" table by ID using prepared statements.
func (r *SnippetsRepo) UpdateSnippet(userId, id int, name, description, content string) (Snippet, error) {
	query := `
        UPDATE snippets
        SET name = ?, description = ?, content = ?
        WHERE id = ? AND user_id = ?
    `
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return Snippet{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, content, id, userId)
	if err != nil {
		log.Println("Error updating snippet:", err)
		return Snippet{}, err
	}

	// Return the updated snippet
	return Snippet{ID: id, Name: name, Description: description, Content: content}, nil
}

// DeleteSnippet deletes a single snippet from the "snippets" table by ID using prepared statements.
func (r *SnippetsRepo) DeleteSnippet(snippetID int) error {
	query := "DELETE FROM snippets WHERE id = ?"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(snippetID)
	if err != nil {
		log.Println("Error deleting snippet:", err)
	}
	return err
}
