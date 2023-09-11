package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type Snippet struct {
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Content     string `json:"content,omitempty"`
}

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}
type Storage struct {
	db   *sql.DB
	Name string
}

func New(dbName string) {
	dbPath := getDbPath(dbName)

	_, err := os.Stat(dbPath)
	if err == nil {
		// File exists, so remove it
		err := os.Remove(dbPath)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("File '%s' has been removed.\n", dbPath)
	}

	file, e := os.Create(dbPath)
	if e != nil {
		log.Fatal(e)
	}

	log.Printf("Db file %s was created\n", dbPath)

	if e := file.Close(); e != nil {
		log.Fatal(e)
	}

}

func GetConnection(dbName string) (*Storage, error) {
	dbConn, err := sql.Open("sqlite3", fmt.Sprintf("./%s", dbName))

	if err != nil {
		log.Fatal(fmt.Sprintf("Something wrong happened when tried get connection to db: %s", dbName))
		return nil, err
	}

	return &Storage{db: dbConn, Name: dbName}, nil
}

func (s *Storage) CloseConnection() {
	err := s.db.Close()
	if err != nil {
		fmt.Errorf("could not close db")
	}
}

func (s *Storage) SeedDb(seedFilePath string) error {
	// Read the SQL from the seed file
	sqlBytes, err := os.ReadFile(seedFilePath)
	if err != nil {
		return err
	}

	sqlQuery := string(sqlBytes)

	_, err = s.db.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetAllSnippets() ([]Snippet, error) {
	// Define the SQL query to select all rows from the "snippets" table
	query := "SELECT * FROM snippets"

	// Execute the query
	rows, err := s.db.Query(query)
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
		if err := rows.Scan(&snippet.ID, &snippet.Name, &snippet.Description, &snippet.Content); err != nil {
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
func (s *Storage) SaveSnippet(name, description, content string) (Snippet, error) {
	query := `
        INSERT INTO snippets (name, description, content)
        VALUES (?, ?, ?)
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return Snippet{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, content)
	if err != nil {
		log.Println("Error saving snippet:", err)
		return Snippet{}, err
	}

	// Retrieve the newly created snippet's ID from the database
	var id int
	err = s.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return Snippet{}, err
	}

	// Return the newly created snippet with the generated ID
	return Snippet{ID: id, Name: name, Description: description, Content: content}, nil
}

// UpdateSnippet updates an existing snippet in the "snippets" table by ID using prepared statements.
func (s *Storage) UpdateSnippet(id int, name, description, content string) (Snippet, error) {
	query := `
        UPDATE snippets
        SET name = ?, description = ?, content = ?
        WHERE id = ?
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return Snippet{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, description, content, id)
	if err != nil {
		log.Println("Error updating snippet:", err)
		return Snippet{}, err
	}

	// Return the updated snippet
	return Snippet{ID: id, Name: name, Description: description, Content: content}, nil
}

// DeleteSnippet deletes a single snippet from the "snippets" table by ID using prepared statements.
func (s *Storage) DeleteSnippet(snippetID int) error {
	query := "DELETE FROM snippets WHERE id = ?"
	stmt, err := s.db.Prepare(query)
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

func getDbPath(dbName string) string {
	return dbName
	//return fmt.Sprintf("./databases/%s", dbName)
}

func SetupNewTestDb(dbName string) {
	New(dbName)

	storage, _ := GetConnection(dbName)

	err := storage.SeedDb("./db/sql/init.sql")

	if err != nil {
		log.Fatal("Could not init db", err)
	}

	err = storage.SeedDb("./db/sql/seed.sql")
	if err != nil {
		log.Fatal("Could not seed db", err)
	}
}

// CreateUser creates a new user and returns the created user.
func (s *Storage) CreateUser(username, email, fullName string) (User, error) {
	query := `
        INSERT INTO users (username, email, full_name)
        VALUES (?, ?, ?)
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return User{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, email, fullName)
	if err != nil {
		log.Println("Error creating user:", err)
		return User{}, err
	}

	// Retrieve the newly created user's ID from the database
	var id int
	err = s.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return User{}, err
	}

	// Retrieve the user's created_at and updated_at timestamps
	var createdAt, updatedAt string
	query = "SELECT created_at, updated_at FROM users WHERE id = ?"
	err = s.db.QueryRow(query, id).Scan(&createdAt, &updatedAt)
	if err != nil {
		log.Println("Error retrieving timestamps:", err)
		return User{}, err
	}

	// Return the newly created user with the generated ID and timestamps
	return User{
		ID:        id,
		Username:  username,
		Email:     email,
		FullName:  fullName,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

// GetUserByID retrieves a user by ID and returns it.
func (s *Storage) GetUserByID(id int) (User, error) {
	query := "SELECT id, username, email, full_name, created_at, updated_at FROM users WHERE id = ?"
	row := s.db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error retrieving user:", err)
		return User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user and returns the updated user.
func (s *Storage) UpdateUser(id int, username, email, fullName string) (User, error) {
	query := `
        UPDATE users
        SET username = ?, email = ?, full_name = ?
        WHERE id = ?
    `
	stmt, err := s.db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return User{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, email, fullName, id)
	if err != nil {
		log.Println("Error updating user:", err)
		return User{}, err
	}

	// Retrieve the updated user's timestamps
	var updatedAt string
	query = "SELECT updated_at FROM users WHERE id = ?"
	err = s.db.QueryRow(query, id).Scan(&updatedAt)
	if err != nil {
		log.Println("Error retrieving timestamps:", err)
		return User{}, err
	}

	// Return the updated user
	return User{
		ID:        id,
		Username:  username,
		Email:     email,
		FullName:  fullName,
		UpdatedAt: updatedAt,
	}, nil
}

// DeleteUser deletes a user by ID.
func (s *Storage) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := s.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting user:", err)
	}
	return err
}
