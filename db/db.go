package db

import (
	"database/sql"
	"log"
	"os"
	"snippetier/db/repo"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db           *sql.DB
	Name         string
	UsersRepo    *repo.UsersRepo
	SnippetsRepo *repo.SnippetsRepo
}

func initRepos(db *sql.DB) *Storage {
	usersRepo := repo.NewUsersRepo(db)
	snippetsRepo := repo.NewSnippetsRepo(db)

	return &Storage{db: db, UsersRepo: usersRepo, SnippetsRepo: snippetsRepo}
}

func GetConnection() (*Storage, error) {
	dbConn, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		return nil, err
	}

	if err = dbConn.Ping(); err != nil {
		return nil, err
	}

	return initRepos(dbConn), nil
}

func (s *Storage) CloseConnection() {
	err := s.db.Close()
	if err != nil {
		log.Fatalf("could not close db")
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
