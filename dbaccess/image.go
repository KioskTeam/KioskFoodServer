package dbaccess

const imagesOfAFoodSql = `
  SELECT i.url image
  FROM images i, foods_images fi
  WHERE fi.image_id = i.id
    AND fi.food_id = $1
`

type imageResult struct {
	string
	error
}

type imageRequest struct {
	id         int64
	resultChan chan<- imageResult
	errChan    chan<- error
}

var imageRequestsChan = make(chan imageRequest)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range imageRequestsChan {
				result := make([]string, 0, 100)

				req.errChan <- db.Select(&result, imagesOfAFoodSql, req.id)

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
