package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (c PostgresConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

func main() {
	c := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", c.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name TEXT,
	email TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS orders (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	amount INT,
	description TEXT
);`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created.")

	name := "Robert Smith"
	email := "singer@thecure.com"

	var id int

	err = db.QueryRow(
		` INSERT INTO users(name, email) VALUES ($1, $2) RETURNING id;`,
		name, email).Scan(&id)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User id '%v' created.\n", id)

	var n, e string
	err = db.QueryRow(
		`SELECT name, email FROM users WHERE id = $1`, id,
	).Scan(&n, &e)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("no rows found!")
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("User information: name=%s email=%s", n, e)
}
