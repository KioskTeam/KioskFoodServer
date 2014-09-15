package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/KioskTeam/KioskFoodServer/dbaccess"
)

func main() {
	http.HandleFunc("/restaurant/", handleRestaurant)

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

func handleRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(r.URL.Path[len("/restaurant/"):])
	if err != nil {
		panic(err)
	}

	data, err := dbaccess.GetRestaurant(int64(restaurantID))
	if err != nil {
		panic(err)
	}

	marshal, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Write(marshal)
}
