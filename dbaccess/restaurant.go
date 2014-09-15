package dbaccess

// Restaurant stores all the data
type Restaurant struct {
	Name       string
	Address    string
	Categories []FoodCategory
}

const restaurantWithIDSql = `
  SELECT r.name, r.address
  FROM restaurants r
  WHERE r.id = $1
`

// GetRestaurant returns all the data of a restaurant with the specified id.
func GetRestaurant(id int64) (Restaurant, error) {
	var (
		restaurant struct {
			Name    string
			Address string
		}

		result Restaurant
	)

	if err := db.Get(&restaurant, restaurantWithIDSql, id); err != nil {
		return result, err
	}

	cats, err := getCatsOfRestaurant(id)
	if err != nil {
		return result, err
	}

	result = Restaurant{
		Name:       restaurant.Name,
		Address:    restaurant.Address,
		Categories: cats,
	}

	return result, nil
}
