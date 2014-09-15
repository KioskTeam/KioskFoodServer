package dbaccess

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

func getFoodsOfCat(id int64) ([]Food, error) {
	var (
		foods []struct {
			ID        int64
			Name      string
			Price     int64
			Thumbnail string
		}

		result = make([]Food, 0, 100)
	)

	if err := db.Select(&foods, foodsOfACatSql, id); err != nil {
		return result, err
	}

	for _, food := range foods {
		pics, err := getPicsOfFood(food.ID)
		if err != nil {
			return result, err
		}

		result = append(result, Food{food.Name, food.Price, food.Thumbnail, pics})
	}

	return result, nil
}
