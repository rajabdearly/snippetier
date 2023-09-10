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

func GetConnection(dbName string) (*sql.DB, error) {
	dbConn, err := sql.Open("sqlite3", fmt.Sprintf("./%s", dbName))

	if err != nil {
		log.Fatal(fmt.Sprintf("Something wrong happened when tried get connection to db: %s", dbName))
		return nil, err
	}

	return dbConn, nil
}

func SeedDb(dbConn *sql.DB, seedFilePath string) error {
	// Read the SQL from the seed file
	sqlBytes, err := os.ReadFile(seedFilePath)
	if err != nil {
		return err
	}

	sqlQuery := string(sqlBytes)

	_, err = dbConn.Exec(sqlQuery)
	if err != nil {
		return err
	}

	return nil
}

func GetAllSnippets(dbConn *sql.DB) ([]Snippet, error) {
	// Define the SQL query to select all rows from the "snippets" table
	query := "SELECT * FROM snippets"

	// Execute the query
	rows, err := dbConn.Query(query)
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

func getDbPath(dbName string) string {
	return dbName
	//return fmt.Sprintf("./databases/%s", dbName)
}
