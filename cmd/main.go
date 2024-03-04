package main

import (
	"fmt"

	"github.com/myrat012/mini-go-memorize/internal/app"
)

func main() {
	// ./assert/database.db
	application, err := app.NewMemorize("database.db")
	if err != nil {
		fmt.Printf("Error Create application.")
		panic(err)
	}

	application.Start()
}
