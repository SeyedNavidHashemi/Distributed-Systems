package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./files"))

	http.Handle("/files/", http.StripPrefix("/files/", fs))

	log.Println("File server running on :8081")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}