-- Create the "users" table with timestamp columns
CREATE TABLE IF NOT EXISTS users (
                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                       username TEXT NOT NULL,
                       email TEXT NOT NULL UNIQUE,
                       created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                       updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger to update the "updated_at" column when a row is updated
CREATE TRIGGER users_update_timestamp
    AFTER UPDATE ON users
BEGIN
    UPDATE users
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

-- Create the "snippets" table with timestamp columns
CREATE TABLE IF NOT EXISTS snippets (
                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                          name TEXT NOT NULL,
                          description TEXT,
                          content TEXT,
                          user_id INTEGER NOT NULL,
                          created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                          updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create a trigger to update the "updated_at" column when a row is updated
CREATE TRIGGER snippets_update_timestamp
    AFTER UPDATE ON snippets
BEGIN
    UPDATE snippets
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

