package main

import (
	"log"
)

func main() {
	r := setupRouter()
	err := r.Run()
	if err != nil {
		log.Fatal(err)
	}
}
