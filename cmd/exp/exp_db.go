package main

import (
	"fmt"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rafaelmdurante/lenslocked/models"
)

func main() {
	c := models.DefaultPostgresConfig()

	// open connection
	db, err := models.Open(c)
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

	// create tables by connecting to the docker
	// docker exec -it lenslocked_db /usr/bin/psql -U <user> -d <dbname>
	// create tables as in models/sql/*.sql files

	// create user
	us := models.UserService{DB: db}
	user, err := us.Create("bob@email.com", "my secret password")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)
}
