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

type resRequest struct {
	id         int64
	resultChan chan<- resResult
	errChan    chan<- error
}

var resRequestsChan = make(chan resRequest)

func init() {
	for i := 0; i < 20; i++ {
		go func() {
			for req := range resRequestsChan {
				var restaurant struct {
					Name       string
					Name_fa    string
					Address    string
					Address_fa string
				}

				req.errChan <- db.Get(&restaurant, restaurantWithIDSql, req.id)

				catsChan, errChan := getCatsOfRestaurant(req.id)

				err := <-errChan
				cats := []FoodCategory{}

				for cat := range catsChan {
					if cat.error == nil {
						cats = append(cats, cat.FoodCategory)
					} else {
						log.Print(cat.error)
					}
				}

				req.resultChan <- resResult{Restaurant{restaurant.Name, restaurant.Name_fa, restaurant.Address, restaurant.Address_fa, cats}, err}

				close(req.resultChan)
			}
		}()
	}
}

// GetRestaurant returns all the data of a restaurant with the specified id.
func GetRestaurant(id int64) (<-chan resResult, <-chan error) {
	c := make(chan resResult)
	errc := make(chan error, 1)

	resRequestsChan <- resRequest{id, c, errc}

	return c, errc
}
