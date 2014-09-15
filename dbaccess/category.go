package dbaccess

import "log"

// FoodCategory is a food category :-)
type FoodCategory struct {
	Name  string
	Image string
	Foods []Food
}

const catsOfARestaurantSql = `
  SELECT c.id, c.name, i.url image
  FROM food_categories c, images i, restaurants_food_categories rc
  WHERE c.image_id = i.id
		AND rc.food_category_id = c.id
		AND rc.restaurant_id = $1
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
				ID    int64
				Name  string
				Image string
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

			c <- catResult{FoodCategory{cat.Name, cat.Image, foods}, err}
		}

		close(c)
	}()
	return c, errc
}
