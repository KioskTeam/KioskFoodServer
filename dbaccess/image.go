package dbaccess

const imagesOfAFoodSql = `
  SELECT i.url image
  FROM images i, foods_images fi
  WHERE fi.image_id = i.id
    AND fi.food_id = $1
`

func getPicsOfFood(id int64) ([]string, error) {
	var (
		pics   []struct{ Image string }
		result = make([]string, 0, 100)
	)

	if err := db.Select(&pics, imagesOfAFoodSql, id); err != nil {
		return result, err
	}

	for _, p := range pics {
		result = append(result, p.Image)
	}

	return result, nil
}
