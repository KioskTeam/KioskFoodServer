package dbaccess

import "log"

// FoodCategory is a food category :-)
type FoodCategory struct {
	Name    string
	Name_fa string
	Image   string
	Foods   []Food
}

const catsOfARestaurantSql = `
  SELECT c.id, c.name, c.name_fa, i.url image
  FROM food_categories c, images i
  WHERE c.image_id = i.id
		AND c.restaurant_id = $1
`

type catResult struct {
	FoodCategory
	error
}

func getCatsOfRestaurant(id int64) (<-chan catResult, <-chan error) {
	c := make(chan catResult)
	errc := make(chan error, 1)

	go func() {
		var (
			cats []struct {
				ID      int64
				Name    string
				Name_fa string
				Image   string
			}
		)

		errc <- db.Select(&cats, catsOfARestaurantSql, id)

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

			c <- catResult{FoodCategory{cat.Name, cat.Name_fa, cat.Image, foods}, err}
		}

		close(c)
	}()
	return c, errc
}
