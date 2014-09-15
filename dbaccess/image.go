package dbaccess

const imagesOfAFoodSql = `
  SELECT i.url image
  FROM images i, foods_images fi
  WHERE fi.image_id = i.id
    AND fi.food_id = $1
`

func getPicsOfFood(id int64) ([]string, error) {
	result := make([]string, 0, 100)

	if err := db.Select(&result, imagesOfAFoodSql, id); err != nil {
		return nil, err
	}

	return result, nil
}
