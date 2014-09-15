package dbaccess

import "log"

// Food stores basic data about foods
type Food struct {
	Name      string
	Price     int64
	Thumbnail string
	Pictures  []string
}

const foodsOfACatSql = `
  SELECT f.id, f.name, f.price, i.url thumbnail
  FROM foods f, images i
  WHERE f.image_id = i.id
    AND f.food_category_id = $1
`

type foodResult struct {
	Food
	error
}

func getFoodsOfCat(id int64) (<-chan foodResult, <-chan error) {
	c := make(chan foodResult)
	errc := make(chan error, 1)

	go func() {

		type dbfood struct {
			ID        int64
			Name      string
			Price     int64
			Thumbnail string
		}

		var foods []dbfood

		errc <- db.Select(&foods, foodsOfACatSql, id)

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

			c <- foodResult{Food{food.Name, food.Price, food.Thumbnail, pics}, err}
		}

		close(c)
	}()

	return c, errc
}
