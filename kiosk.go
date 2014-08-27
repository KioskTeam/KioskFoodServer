package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", root)

	port := getPort()
	fmt.Println("listenning on", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		fmt.Println("Warning: $PORT is not set!")
	}
	return ":" + port
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello world!")
}
