package dbaccess

import (
	"log"
	"sync"
	"time"
)

const restaurantWithIDSql = `
  SELECT r.name, r.name_fa, r.address, r.address_fa
  FROM restaurants r
  WHERE r.id = $1
`

type (
	// Restaurant stores all the data
	Restaurant struct {
		Name       string
		Name_fa    string
		Address    string
		Address_fa string
		Categories []FoodCategory
		Galery     []GaleryImage
	}

	resResult struct {
		Restaurant
		Error error
	}

	resRequest struct {
		id         int64
		resultChan chan<- resResult
		errChan    chan<- error
	}

	dbrestaurant struct {
		Name       string
		Name_fa    string
		Address    string
		Address_fa string
	}

	resCacheStore struct {
		time time.Time
		res  dbrestaurant
	}
)

var (
	resRequestsChan = make(chan resRequest)
	resCache        = struct {
		sync.RWMutex
		c map[int64]resCacheStore
	}{c: make(map[int64]resCacheStore)}
)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range resRequestsChan {
				var restaurant dbrestaurant

				resCache.RLock()
				cache, ok := resCache.c[req.id]
				resCache.RUnlock()

				if ok && CacheIsRecent(cache.time) {
					req.errChan <- nil
					restaurant = cache.res
				} else {
					req.errChan <- db.Get(&restaurant, restaurantWithIDSql, req.id)
				}

				catsChan, errChan := getCatsOfRestaurant(req.id)

				err := <-errChan
				cats := []FoodCategory{}

				for cat := range catsChan {
					if cat.error == nil {
						cats = append(cats, cat.FoodCategory)
					} else {
						log.Print(cat.error)
					}
				}

				galeryChan, errChan := getGaleryOfRestaurant(req.id)

				err = <-errChan
				galeryImages := []GaleryImage{}

				for galImg := range galeryChan {
					if galImg.error == nil {
						galeryImages = append(galeryImages, galImg.GaleryImage)
					} else {
						log.Print(galImg.error)
					}
				}

				resie := Restaurant{
					restaurant.Name, restaurant.Name_fa,
					restaurant.Address, restaurant.Address_fa,
					cats, galeryImages,
				}
				req.resultChan <- resResult{resie, err}

				resCache.Lock()
				resCache.c[req.id] = resCacheStore{time.Now(), restaurant}
				resCache.Unlock()

				close(req.resultChan)
			}
		}()
	}
}

// GetRestaurant returns all the data of a restaurant with the specified id.
func GetRestaurant(id int64) (<-chan resResult, <-chan error) {
	c := make(chan resResult)
	errc := make(chan error, 1)

	resRequestsChan <- resRequest{id, c, errc}

	return c, errc
}
