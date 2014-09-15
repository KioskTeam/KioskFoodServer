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

func getPicsOfFood(id int64) (<-chan imageResult, <-chan error) {
	c := make(chan imageResult)
	errc := make(chan error, 1)

	go func() {
		result := make([]string, 0, 100)

		errc <- db.Select(&result, imagesOfAFoodSql, id)

		for i := 0; i < len(result); i++ {
			c <- imageResult{result[i], nil}
		}

		close(c)
	}()

	return c, errc
}
