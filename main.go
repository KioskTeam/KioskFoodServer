package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

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

type responseCacheStore struct {
	time time.Time
	resp []byte
}

var responseCache = struct {
	sync.RWMutex
	c map[int]responseCacheStore
}{c: make(map[int]responseCacheStore)}

func handleRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.Atoi(r.URL.Path[len("/restaurant/"):])
	if err != nil {
		panic(err)
	}

	responseCache.RLock()
	cache, ok := responseCache.c[restaurantID]
	responseCache.RUnlock()

	var response []byte

	if ok && dbaccess.CacheIsRecent(cache.time) {
		response = cache.resp
	} else {
		dataChan, errChan := dbaccess.GetRestaurant(int64(restaurantID))
		if err := <-errChan; err != nil {
			panic(err)
		}

		data := <-dataChan
		if data.Error != nil {
			panic(data.Error)
		}

		var err error
		response, err = json.Marshal(data.Restaurant)
		if err != nil {
			panic(err)
		}

		responseCache.Lock()
		responseCache.c[restaurantID] = responseCacheStore{time.Now(), response}
		responseCache.Unlock()
	}

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(response)
}
