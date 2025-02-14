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
	db, err := sql.Open("postgres", "host=127.0.0.1 port=5432 user=postgres password=12345 dbname=mydatabase sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	users, err := getusers(db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(users)

	user, err := getUserByID(db, 1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

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
	// 	log.Fatal("Ошибка при создании таблицы: ", err)
	// }

	// insertQuery := `
	// INSERT INTO users (name, email, password, registered_at)
	// VALUES
	// 	('John Doe', 'john.doe@example.com', 'password123', NOW()),
	// 	('Jane Smith', 'jane.smith@example.com', 'password456', NOW()),
	// 	('Alice Johnson', 'alice.johnson@example.com', 'password789', NOW())
	// `

	// _, err = db.Exec(insertQuery)
	// if err != nil {
	// 	log.Fatal("Ошибка при вставке данных: ", err)
	// }

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

func getUserByID(db *sql.DB, id int64) (User, error) {
	var u User
	err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).
		Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
	if err != nil {
		return User{}, err
	}
	return u, nil
}
