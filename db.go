package main

import (
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
	cats, err := getAllCats()
	if err != nil {
		return Restaurant{}, err
	}

	restaurant := Restaurant{
		Name:       "Good Father",
		Address:    "Tehran",
		Categories: cats,
	}

	return restaurant, nil
}

func getAllCats() ([]FoodCategory, error) {
	var (
		cats []struct {
			Id    int64
			Name  string
			Image string
		}

		result []FoodCategory
	)

	if err := db.Select(&cats, allCatsSql); err != nil {
		return result, err
	}

	for _, cat := range cats {
		foods, err := getFoodsOfCat(cat.Id)
		if err != nil {
			return result, err
		}

		result = append(result, FoodCategory{cat.Name, cat.Image, foods})
	}

	return result, nil
}

func getFoodsOfCat(id int64) ([]Food, error) {
	var (
		foods []struct {
			Id        int64
			Name      string
			Price     int64
			Thumbnail string
		}

		result []Food
	)

	if err := db.Select(&foods, foodsOfACatSql, id); err != nil {
		return result, err
	}

	for _, food := range foods {
		pics, err := getPicsOfFood(food.Id)
		if err != nil {
			return result, err
		}

		result = append(result, Food{food.Name, food.Price, food.Thumbnail, pics})
	}

	return result, nil
}

func getPicsOfFood(id int64) ([]string, error) {
	var (
		pics   []struct{ Image string }
		result []string
	)

	if err := db.Select(&pics, imagesOfAFoodSql, id); err != nil {
		return result, err
	}

	for _, p := range pics {
		result = append(result, p.Image)
	}

	return result, nil
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
	Price     int64
	Thumbnail string
	Pictures  []string
}
