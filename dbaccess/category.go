package dbaccess

import (
	"log"
	"sync"
	"time"
)

const catsOfARestaurantSql = `
  SELECT c.id, c.name, c.name_fa, i.url image
  FROM food_categories c, images i
  WHERE c.image_id = i.id
		AND c.restaurant_id = $1
`

type (
	// FoodCategory is a food category :-)
	FoodCategory struct {
		Name    string
		Name_fa string
		Image   string
		Foods   []Food
	}

	catResult struct {
		FoodCategory
		error
	}

	catRequest struct {
		id         int64
		resultChan chan<- catResult
		errChan    chan<- error
	}

	catCacheStore struct {
		time time.Time
		cats []FoodCategory
	}
)

var (
	catRequestsChan = make(chan catRequest)
	catCache        = struct {
		sync.RWMutex
		c map[int64]catCacheStore
	}{c: make(map[int64]catCacheStore)}
)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range catRequestsChan {
				catCache.RLock()
				cache, ok := catCache.c[req.id]
				catCache.RUnlock()

				if ok && CacheIsRecent(cache.time) {
					req.errChan <- nil
					for _, cat := range cache.cats {
						req.resultChan <- catResult{cat, nil}
					}
				} else {
					allCaties := make([]FoodCategory, 0, 10)

					var cats []struct {
						ID      int64
						Name    string
						Name_fa string
						Image   string
					}

					req.errChan <- db.Select(&cats, catsOfARestaurantSql, req.id)

					for _, cat := range cats {
						foodsChan, errChan := getFoodsOfCat(cat.ID)

						err := <-errChan
						foods := []Food{}

						for food := range foodsChan {
							if food.error == nil {
								foods = append(foods, food.Food)
							} else {
								log.Print(food.error)
							}
						}

						catie := FoodCategory{cat.Name, cat.Name_fa, cat.Image, foods}
						allCaties = append(allCaties, catie)

						req.resultChan <- catResult{catie, err}
					}

					catCache.Lock()
					catCache.c[req.id] = catCacheStore{time.Now(), allCaties}
					catCache.Unlock()
				}

				close(req.resultChan)
			}
		}()
	}
}

func getCatsOfRestaurant(id int64) (<-chan catResult, <-chan error) {
	c := make(chan catResult)
	errc := make(chan error, 1)

	catRequestsChan <- catRequest{id, c, errc}

	return c, errc
}
