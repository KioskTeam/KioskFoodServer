package dbaccess

import "log"

// Restaurant stores all the data
type Restaurant struct {
	Name       string
	Name_fa    string
	Address    string
	Address_fa string
	Categories []FoodCategory
}

const restaurantWithIDSql = `
  SELECT r.name, r.name_fa, r.address, r.address_fa
  FROM restaurants r
  WHERE r.id = $1
`

type resResult struct {
	Restaurant
	Error error
}

// GetRestaurant returns all the data of a restaurant with the specified id.
func GetRestaurant(id int64) (<-chan resResult, <-chan error) {
	c := make(chan resResult)
	errc := make(chan error, 1)

	go func() {
		var (
			restaurant struct {
				Name       string
				Name_fa    string
				Address    string
				Address_fa string
			}
		)

		errc <- db.Get(&restaurant, restaurantWithIDSql, id)

		catsChan, errChan := getCatsOfRestaurant(id)

		err := <-errChan
		cats := []FoodCategory{}

		for cat := range catsChan {
			if cat.error == nil {
				cats = append(cats, cat.FoodCategory)
			} else {
				log.Print(cat.error)
			}
		}

		c <- resResult{Restaurant{restaurant.Name, restaurant.Name_fa, restaurant.Address, restaurant.Address_fa, cats}, err}

		close(c)
	}()

	return c, errc
}
