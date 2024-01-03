package main

import (
	"html/template"
	"os"
)

func main() {
	// this is relative to where `go run` runs!
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	user := struct {
		Name string
	}{
		"Raf",
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
