
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO images (id, url)
VALUES
    (1, '/assets/IMG_0031.JPG'),
    (2, '/assets/IMG_0032.JPG'),
    (3, '/assets/IMG_0033.JPG'),
    (4, '/assets/IMG_0034.JPG');

INSERT INTO restaurants (id, name, address)
VALUES
    (1, 'Good Father', 'Tehran');

INSERT INTO food_categories (id, name, image_id)
VALUES
    (1, 'cat 1', 1),
    (2, 'cat 2', 2),
    (3, 'cat 3', 3);

INSERT INTO restaurants_food_categories (id, restaurant_id, food_category_id)
VALUES
    (1, 1, 1),
    (2, 1, 2),
    (3, 1, 3);

INSERT INTO foods (id, food_category_id, name, description, price, image_id)
VALUES
    (1, 1, 'sample food 1', 'some description', 10000, 1),
    (2, 1, 'sample food 2', 'some description', 12000, 1),
    (3, 1, 'sample food 3', 'some description', 15000, 2),
    (4, 2, 'sample food 4', 'some description', 1000, 3),
    (5, 2, 'sample food 5', NULL, 1500, 4);

INSERT INTO foods_images (food_id, image_id)
VALUES
    (1, 1),
    (1, 2),
    (1, 3),
    (2, 4),
    (3, 4);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE FROM foods_images;
DELETE FROM foods;
DELETE FROM restaurants_food_categories;
DELETE FROM food_categories;
DELETE FROM restaurants;
DELETE FROM images;
