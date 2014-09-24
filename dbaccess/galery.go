package dbaccess

import (
	"sync"
	"time"
)

const galeryOfARestaurant = `
  SELECT g.description, i.url image
  FROM images i, restaurant_galery g
  WHERE g.image_id = i.id
    AND g.restaurant_id = $1
`

type (
	GaleryImage struct {
		Image       string
		Description string
	}

	galeryResult struct {
		GaleryImage
		error
	}

	galeryRequest struct {
		id         int64
		resultChan chan<- galeryResult
		errChan    chan<- error
	}

	galeryCacheStore struct {
		time   time.Time
		galery []GaleryImage
	}
)

var (
	galeryRequestsChan = make(chan galeryRequest)
	galeryCache        = struct {
		sync.RWMutex
		c map[int64]galeryCacheStore
	}{c: make(map[int64]galeryCacheStore)}
)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range galeryRequestsChan {
				var result []GaleryImage

				galeryCache.RLock()
				cache, ok := galeryCache.c[req.id]
				galeryCache.RUnlock()

				if ok && CacheIsRecent(cache.time) {
					req.errChan <- nil
					result = cache.galery
				} else {
					req.errChan <- db.Select(&result, galeryOfARestaurant, req.id)
					galeryCache.Lock()
					galeryCache.c[req.id] = galeryCacheStore{time.Now(), result}
					galeryCache.Unlock()
				}

				for _, img := range result {
					req.resultChan <- galeryResult{img, nil}
				}

				close(req.resultChan)
			}
		}()
	}
}

func getGaleryOfRestaurant(id int64) (<-chan galeryResult, <-chan error) {
	c := make(chan galeryResult)
	errc := make(chan error, 1)

	galeryRequestsChan <- galeryRequest{id, c, errc}

	return c, errc
}
