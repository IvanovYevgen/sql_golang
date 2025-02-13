package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
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


	insertQuery := `
	INSERT INTO users (name, email, password, registered_at) 
	VALUES
		('John Doe', 'john.doe@example.com', 'password123', NOW()),
		('Jane Smith', 'jane.smith@example.com', 'password456', NOW()),
		('Alice Johnson', 'alice.johnson@example.com', 'password789', NOW())
	`

	_, err = db.Exec(insertQuery)
	if err != nil {
		log.Fatal("Ошибка при вставке данных: ", err)
	}


	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		u := User{}
		err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.RegisteredAt)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(users)
}
