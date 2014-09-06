package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/api/latest", latest)

	assetsDir := http.Dir(os.Getenv("STATIC_ASSETS_PATH"))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(assetsDir)))

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

func latest(w http.ResponseWriter, r *http.Request) {
	data, err := getAllData()
	if err != nil {
		panic(err)
	}
	marshal, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write(marshal)
}
