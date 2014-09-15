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

func getCatsOfRestaurant(id int64) ([]FoodCategory, error) {
	var (
		cats []struct {
			ID    int64
			Name  string
			Image string
		}

		result = make([]FoodCategory, 0, 100)
	)

	if err := db.Select(&cats, catsOfARestaurantSql, id); err != nil {
		return result, err
	}

	for _, cat := range cats {
		foodsChan, errChan := getFoodsOfCat(cat.ID)
		if err := <-errChan; err != nil {
			return result, err
		}

		foods := []Food{}
		for food := range foodsChan {
			if food.error == nil {
				foods = append(foods, food.Food)
			} else {
				log.Print(food.error)
			}
		}

		result = append(result, FoodCategory{cat.Name, cat.Image, foods})
	}

	return result, nil
}
