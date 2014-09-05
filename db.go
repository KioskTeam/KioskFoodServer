package main

import (
	"database/sql"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db = sqlx.MustConnect("postgres", os.Getenv("DATABASE_URL"))

	allCatsSql = `
    SELECT c.id, c.name, i.url image
    FROM food_categories c, images i
    WHERE c.image_id = i.id
  `

	foodsOfACatSql = `
    SELECT f.id, f.name, f.price, i.url thumbnail
    FROM foods f, images i
    WHERE f.image_id = i.id
      AND f.food_category_id = $1
  `

	imagesOfAFoodSql = `
    SELECT i.url image
    FROM images i, foods_images fi
    WHERE fi.image_id = i.id
      AND fi.food_id = $1
  `
)

func getAllData() (Restaurant, error) {
	rows, err := db.Query(allCatsSql)
	if err != nil {
		return Restaurant{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id    int64
			name  string
			image string
		)

		err = rows.Scan(&id, &name, &image)
		if err != nil {
			return Restaurant{}, err
		}

		cat := FoodCategory{name, image, make([]Food, 0, 100)}

		rows2, err := db.Query(foodsOfACatSql, id)
		if err != nil {
			return Restaurant{}, err
		}
		defer rows2.Close()

		for rows2.Next() {
			var (
				id        int64
				name      string
				price     sql.NullInt64
				thumbnail string
			)

			err = rows2.Scan(&id, &name, &price, &thumbnail)
			if err != nil {
				return Restaurant{}, err
			}

			food := Food{name, price, thumbnail, make([]string, 0, 100)}

			var pics []struct{ Image string }
			if err = db.Select(&pics, imagesOfAFoodSql, id); err != nil {
				return Restaurant{}, err
			}
			for _, p := range pics {
				food.Pictures = append(food.Pictures, p.Image)
			}
			cat.Foods = append(cat.Foods, food)
		}

		restaurant.Categories = append(restaurant.Categories, cat)
	}

	return restaurant, nil
}

// Restaurant stores all the data
type Restaurant struct {
	Name       string
	Address    string
	Categories []FoodCategory
}

// FoodCategory is a food category :-)
type FoodCategory struct {
	Name  string
	Image string
	Foods []Food
}

// Food stores basic data about foods
type Food struct {
	Name      string
	Price     sql.NullInt64
	Thumbnail string
	Pictures  []string
}

// Data is a temp stub for database
var restaurant = Restaurant{
	Name:       "Good Father",
	Address:    "Tehran",
	Categories: make([]FoodCategory, 0, 100),
}
