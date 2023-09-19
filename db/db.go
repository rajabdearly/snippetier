package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"snippetier/db/repo"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db           *sql.DB
	Name         string
	UsersRepo    *repo.UsersRepo
	SnippetsRepo *repo.SnippetsRepo
}

func createDbtest(dbName string) error {
	dbPath := getDbPath(dbName)

	_, err := os.Stat(dbPath)
	if err == nil {
		// File exists, so remove it
		err := os.Remove(dbPath)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		fmt.Printf("File '%s' has been removed.\n", dbPath)
	} else if os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return err
		}

		log.Printf("Db file %s was created\n", dbPath)

		if err := file.Close(); err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

func initRepos(db *sql.DB) *Storage {
	usersRepo := repo.NewUsersRepo(db)
	snippetsRepo := repo.NewSnippetsRepo(db)

	return &Storage{db: db, UsersRepo: usersRepo, SnippetsRepo: snippetsRepo}
}

func GetConnection(dbName string) (*Storage, error) {
	dbConn, err := sql.Open("sqlite3", fmt.Sprintf("./%s", dbName))

	if err != nil {
		log.Fatalf("Something wrong happened when tried get connection to db: %s", dbName)
		return nil, err
	}

	return initRepos(dbConn), nil
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

func getDbPath(dbName string) string {
	return dbName
	//return fmt.Sprintf("./databases/%s", dbName)
}

func SetupNewTestDb(dbName string) {
	createDbtest(dbName)

	storage, _ := GetConnection(dbName)

	err := storage.SeedDb("./db/sql/init.sql")

	if err != nil {
		log.Fatal("Could not init db:\n", err)
	}

	err = storage.SeedDb("./db/sql/seed.sql")
	if err != nil {
		log.Fatal("Could not seed db", err)
	}
}
