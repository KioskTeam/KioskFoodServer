package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

// Restaurant stores all the data
type Restaurant struct {
	Name       string
	Address    string
	Categories []FoodCategory
}

// FoodCategory is a food category :-)
type FoodCategory struct {
	Name  string
	Image string
	Foods []Food
}

// Food stores basic data about foods
type Food struct {
	Name      string
	Price     int
	Thumbnail string
	Pictures  []string
}

// Data is a temp stub for database
var restaurant = Restaurant{
	Name:    "Good Father",
	Address: "Tehran",
	Categories: []FoodCategory{
		{
			Name:  "Pizza",
			Image: "/pizza.jpg",
			Foods: []Food{
				{
					Name:      "Peperony",
					Price:     10000,
					Thumbnail: "/peperony-small.jpg",
					Pictures: []string{
						"/peperony1.jpg",
						"/peperony2.jpg",
					},
				},
			},
		},
	},
}

func latest(w http.ResponseWriter, r *http.Request) {
	marshal, err := json.Marshal(restaurant)
	if err != nil {
		panic(err)
	}
	w.Write(marshal)
}
