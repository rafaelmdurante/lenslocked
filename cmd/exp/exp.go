package main

import (
	"html/template"
	"os"
)

type User struct {
	Name             string
	Age              int
	Address          Address
	Salary           float32
	CountriesVisited []string
	LanguagesFluency map[string]string
}

type Address struct {
	Street string
	City   string
}

func main() {
	// this is relative to where `go run` runs!
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Raf",
		Age:  39,
		Address: Address{
			Street: "Leona Drive",
			City:   "Dun Laoghaire",
		},
		Salary:           1234.56,
		CountriesVisited: []string{"Japan", "Italy", "Scotland"},
		LanguagesFluency: map[string]string{
			"Portuguese": "Native",
			"English":    "Fluent",
			"Japanese":   "Intermediate",
		},
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
