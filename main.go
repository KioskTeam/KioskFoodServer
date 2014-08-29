package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	http.HandleFunc("/api/latest", latest)

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

// Food stores basic data about foods
type Food struct {
	Name  string
	Desc  string
	Price int
	Image string
}

// DataList gets json-serialized to clients
type DataList struct {
	Checksome string
	Foods     []Food
}

// Data is a temp stub for database
var Data = DataList{
	Checksome: strconv.FormatInt(time.Now().Unix(), 10),
	Foods: []Food{
		{"Pasta", "It's delicious", 1000, "/110.jpg"},
	},
}

func latest(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(Data)
	if err != nil {
		panic(err)
	}
	w.Write(marshal)
}
