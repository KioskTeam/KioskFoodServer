package dbaccess

import (
	"sync"
	"time"
)

const imagesOfAFoodSql = `
  SELECT i.url image
  FROM images i, foods_images fi
  WHERE fi.image_id = i.id
    AND fi.food_id = $1
`

type (
	imageResult struct {
		string
		error
	}

	imageRequest struct {
		id         int64
		resultChan chan<- imageResult
		errChan    chan<- error
	}

	imageCacheStore struct {
		time   time.Time
		images []string
	}
)

var (
	imageRequestsChan = make(chan imageRequest)
	imageCache        = struct {
		sync.RWMutex
		c map[int64]imageCacheStore
	}{c: make(map[int64]imageCacheStore)}
)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range imageRequestsChan {
				var result []string

				imageCache.RLock()
				cache, ok := imageCache.c[req.id]
				imageCache.RUnlock()

				if ok && cacheIsRecent(cache.time) {
					req.errChan <- nil
					result = cache.images
				} else {
					req.errChan <- db.Select(&result, imagesOfAFoodSql, req.id)
					imageCache.Lock()
					imageCache.c[req.id] = imageCacheStore{time.Now(), result}
					imageCache.Unlock()
				}

				for _, img := range result {
					req.resultChan <- imageResult{img, nil}
				}

				close(req.resultChan)
			}
		}()
	}
}

func getPicsOfFood(id int64) (<-chan imageResult, <-chan error) {
	r, e := make(chan imageResult), make(chan error, 1)
	imageRequestsChan <- imageRequest{id, r, e}
	return r, e
}
