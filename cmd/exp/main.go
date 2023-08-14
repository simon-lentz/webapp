package main

import (
	"fmt"
	"os"

	"html/template"
)

func main() {
	fmt.Println("Experimental main.go")
	t, err := template.ParseFiles("hello.html")
	if err != nil {
		panic(err)
	}

	user := struct {
		Name string
	}{
		Name: "Simon",
	}

	if err = t.Execute(os.Stdout, user); err != nil {
		panic(err)
	}
}
