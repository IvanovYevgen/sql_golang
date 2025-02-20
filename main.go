package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type User struct {
	ID           int64
	Name         string
	Email        string
	Password     string
	RegisteredAt time.Time
}

func main() {
	// // Delete the database
	// if err := dropDatabase("mydatabase"); err != nil {
	// 	log.Fatal("Error in deleting DB:", err)
	// } else {
	// 	fmt.Println("DB deleted")
	// }

	// Create the database before starting work
	// if err := createDatabase("mydatabase"); err != nil {
	// 	log.Fatal("Error creating database:", err)
	// }

	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=12345 dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// 	createLogsTableQuery := `
	// CREATE TABLE IF NOT EXISTS logs (
	// 	id SERIAL PRIMARY KEY,
	// 	entity VARCHAR(100),
	// 	action VARCHAR(100),
	// 	timestamp TIMESTAMP DEFAULT NOW()
	// );`

	// 	_, err = db.Exec(createLogsTableQuery)
	// 	if err != nil {
	// 		log.Fatal("Error creating logs table: ", err)
	// 	}

	// createTableQuery := `
	// CREATE TABLE IF NOT EXISTS users (
	// 	id SERIAL PRIMARY KEY,
	// 	name VARCHAR(100),
	// 	email VARCHAR(100) UNIQUE,
	// 	password VARCHAR(100),
	// 	registered_at TIMESTAMP
	// );`

	// _, err = db.Exec(createTableQuery)
	// if err != nil {
	// 	log.Fatal("Error with creating a table ", err)
	// }

	// Inserting data
	err = insertUser(db, User{
		Name:  "Huikol lOolll",
		Email: "huilik.com", Password: "passworsdfasdf"})
	if err != nil {
		log.Fatal(err)
	}

	// insertQuery := `
	// INSERT INTO users (name, email, password, registered_at)
	// VALUES
	// 	('John Doe', 'john.doe@example.com', 'password123', NOW()),
	// 	('Jane Smith', 'jane.smith@example.com', 'password456', NOW()),
	// 	('Alice Johnson', 'alice.johnson@example.com', 'password789', NOW())
	// `

	// _, err = db.Exec(insertQuery)
	// if err != nil {
	// 	log.Fatal("Error with inserting data ", err)
	// }

	// users, err := getusers(db)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(users)

	// error := deleteUser(db, 1)
	// if (error) != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("User deleted")
	// }

	// error := updateUser(db, 2, User{
	// 	Name:  "Jana D'ark",
	// 	Email: "janna@sdf.co", Password: "password123", RegisteredAt: time.Now(),
	// })
	// if (error) != nil {
	// 	log.Fatal(err)
	// } else {
	// 	fmt.Println("User updated")
	// }

	users, err := getusers(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)

	// user, err := getUserByID(db, 1)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(user)

}

func getusers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// func getUserByID(db *sql.DB, id int64) (User, error) {
// 	var u User
// 	err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).
// 		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
// 	if err != nil {
// 		return User{}, err
// 	}
// 	return u, nil
// }

func insertUser(db *sql.DB, u User) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO users (name, email, password, registered_at) VALUES ($1, $2, $3, NOW())",
		u.Name, u.Email, u.Password)

	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO logs (entity, action) VALUES ($1, $2)",
		"users", "created")
	if err != nil {
		return err
	}
	return tx.Commit()
}

// dropDatabase deletes the specified database.
func dropDatabase(dbName string) error {
	// Get a connection to the database
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=12345 dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	// Close all connections to the database
	_, err = db.Exec(fmt.Sprintf(`
		SELECT pg_terminate_backend(pg_stat_activity.pid)
		FROM pg_stat_activity
		WHERE pg_stat_activity.datname = '%s'
		AND pid <> pg_backend_pid();
	`, dbName))
	if err != nil {
		return err
	}

	// Delete database
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s;", dbName))
	if err != nil {
		return err
	}

	return nil
}

func createDatabase(dbName string) error {
	// Connect to the "postgres" database, as the target database does not exist yet
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=12345 dbname=postgres sslmode=disable")
	if err != nil {
		return err
	}
	defer db.Close()

	// Check if the database already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		fmt.Println("Database already exists.")
		return nil
	}

	// Create the new database
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName))
	if err != nil {
		return err
	}

	fmt.Println("Database created successfully:", dbName)
	return nil
}

// Delete user by id
func deleteUser(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	return err
}

// Update user by id
func updateUser(db *sql.DB, id int64, u User) error {
	_, err := db.Exec("UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4",
		u.Name, u.Email, u.Password, id)
	return err
}
