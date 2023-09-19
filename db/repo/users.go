package repo

import (
	"database/sql"
	"log"
)

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"fullName"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{db}
}

// CreateUser creates a new user and returns the created user.
func (r *UsersRepo) CreateUser(username, email, fullName string) (User, error) {
	query := `
        INSERT INTO users (username, email, full_name)
        VALUES (?, ?, ?)
    `
	stmt, err := r.db.Prepare(query)
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

	// Retrieve the newly created user'r ID from the database
	var id int
	err = r.db.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	if err != nil {
		log.Println("Error getting last insert ID:", err)
		return User{}, err
	}

	// Retrieve the user'r created_at and updated_at timestamps
	var createdAt, updatedAt string
	query = "SELECT created_at, updated_at FROM users WHERE id = ?"
	err = r.db.QueryRow(query, id).Scan(&createdAt, &updatedAt)
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
func (r *UsersRepo) GetUserByID(id int) (User, error) {
	query := "SELECT id, username, email, full_name, created_at, updated_at FROM users WHERE id = ?"
	row := r.db.QueryRow(query, id)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error retrieving user:", err)
		return User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user and returns the updated user.
func (r *UsersRepo) UpdateUser(id int, username, email, fullName string) (User, error) {
	query := `
        UPDATE users
        SET username = ?, email = ?, full_name = ?
        WHERE id = ?
    `
	stmt, err := r.db.Prepare(query)
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

	// Retrieve the updated user'r timestamps
	var updatedAt string
	query = "SELECT updated_at FROM users WHERE id = ?"
	err = r.db.QueryRow(query, id).Scan(&updatedAt)
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
func (r *UsersRepo) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Println("Error deleting user:", err)
	}
	return err
}
