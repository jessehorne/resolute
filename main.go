package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jessehorne/resolute/resolute"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/v1", resolute.ServerHandler)
	fmt.Println("starting on port ", os.Getenv("APP_PORT"))
	if err := http.ListenAndServe(":"+os.Getenv("APP_PORT"), nil); err != nil {
		log.Fatalln(err)
	}
}
