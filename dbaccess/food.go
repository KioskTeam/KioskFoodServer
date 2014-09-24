package dbaccess

import "log"

// Food stores basic data about foods
type Food struct {
	Name           string
	Name_fa        string
	Description    string
	Description_fa string
	Price          int64
	Thumbnail      string
	Pictures       []string
}

const foodsOfACatSql = `
  SELECT f.id, f.name, f.name_fa, f.description, f.description_fa, f.price,
		i.url thumbnail
  FROM foods f, images i
  WHERE f.image_id = i.id
    AND f.food_category_id = $1
`

type foodResult struct {
	Food
	error
}

type foodRequest struct {
	id         int64
	resultChan chan<- foodResult
	errChan    chan<- error
}

var foodRequestsChan = make(chan foodRequest)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range foodRequestsChan {

				type dbfood struct {
					ID             int64
					Name           string
					Name_fa        string
					Description    string
					Description_fa string
					Price          int64
					Thumbnail      string
				}

				var foods []dbfood

				req.errChan <- db.Select(&foods, foodsOfACatSql, req.id)

				for _, food := range foods {
					picsChan, picErrc := getPicsOfFood(food.ID)

					err := <-picErrc
					pics := []string{}

					for pic := range picsChan {
						if pic.error == nil {
							pics = append(pics, pic.string)
						} else {
							log.Print(pic.error)
						}
					}

					req.resultChan <- foodResult{Food{food.Name, food.Name_fa, food.Description, food.Description_fa, food.Price, food.Thumbnail, pics}, err}
				}

				close(req.resultChan)
			}
		}()
	}
}

func getFoodsOfCat(id int64) (<-chan foodResult, <-chan error) {
	c := make(chan foodResult)
	errc := make(chan error, 1)

	foodRequestsChan <- foodRequest{id, c, errc}

	return c, errc
}
