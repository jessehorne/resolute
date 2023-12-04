package main

import (
	"fmt"
	"os"

	"github.com/jessehorne/resolute/pkg/v1/resolute"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	host := fmt.Sprintf("%s:%s", "127.0.0.1", os.Getenv("APP_PORT"))
	fmt.Println("listening on:", host)
	s := resolute.NewServer("/v1", host)
	s.Listen()
}
