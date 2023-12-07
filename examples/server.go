package main

import (
	"fmt"

	"github.com/jessehorne/resolute/pkg/v1/resolute"
)

func main() {
	host := "127.0.0.1:5656"

	fmt.Println("listening on:", host)

	s := resolute.NewServer("/v1", host)
	if err := s.Listen("./cert.pem", "./key.pem"); err != nil {
		fmt.Println(err)
	}
}
