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

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
}

func (c PostgresConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Database, c.SSLMode)
}

func createTables(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
	userID SERIAL PRIMARY KEY,
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
}

func createUser(db *sql.DB, name, email string, userID *int) {
	err := db.QueryRow(
		` INSERT INTO users(name, email) VALUES ($1, $2) RETURNING userID;`,
		name, email).Scan(userID)

	if err != nil {
		panic(err)
	}

	fmt.Printf("User userID '%v' created.\n", *userID)
}

func createFakeOrders(db *sql.DB, userID int, quantity int) {
	for i := 1; i <= quantity; i++ {
		amount := i * 100
		desc := fmt.Sprintf("fake order #%d", i)
		_, err := db.Exec(`
INSERT INTO orders(user_id, amount, description) VALUES ($1, $2, $3)`,
			userID, amount, desc)
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Created fake orders.")
}

func findUserByID(db *sql.DB, userID int) {
	var n, e string
	err := db.QueryRow(
		`SELECT name, email FROM users WHERE userID = $1`, userID,
	).Scan(&n, &e)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println("no rows found!")
	}
	if err != nil {
		panic(err)
	}
	fmt.Printf("User information: name=%s email=%s\n", n, e)
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

	// open connection
	db, err := sql.Open("pgx", c.String())
	if err != nil {
		panic(err)
	}
	// ensure connection will be closed when main function finishes
	defer db.Close()

	// test connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	// create tables
	createTables(db)

	// create user
	var userID int
	createUser(db, "Robert Smith", "singer@thecure.com", &userID)

	fmt.Println(userID)
	// select single user
	findUserByID(db, userID)

	// create fake orders
	createFakeOrders(db, userID, 5)

	// read orders
	findAllOrders(db, userID)
}

func findAllOrders(db *sql.DB, userID int) {
	// find all rows
	rows, err := db.Query(
		`SELECT id ,amount, description FROM orders WHERE user_id =$1`,
		userID)

	// check for errors
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// add all rows to slice of order
	var orders []Order

	// sql.Rows was designed to work this way, therefore we are not skipping
	// the first record if we call .Next() immediately
	for rows.Next() {
		var order Order

		order.UserID = userID

		// store a single row data into an object
		err := rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}

		// append the order to the list
		orders = append(orders, order)
	}

	// check for errors
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Orders found:", orders)
}
