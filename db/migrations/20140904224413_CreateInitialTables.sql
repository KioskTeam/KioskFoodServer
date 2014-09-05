
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE images (
  id  serial PRIMARY KEY,
  url varchar(100) NOT NULL,

  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp,
  deleted_at timestamp DEFAULT null
);

CREATE TABLE food_categories (
  id       serial PRIMARY KEY,
  name     varchar(40) NOT NULL,
  image_id integer NOT NULL REFERENCES images,

  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp,
  deleted_at timestamp DEFAULT null
);

CREATE TABLE foods (
  id               serial PRIMARY KEY,
  food_category_id integer NOT NULL REFERENCES food_categories,
  name             varchar(40) NOT NULL,
  description      varchar(200),
  price            integer,
  image_id         integer NOT NULL REFERENCES images,

  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp,
  deleted_at timestamp DEFAULT null
);

CREATE TABLE foods_images (
  id       serial PRIMARY KEY,
  food_id  integer NOT NULL REFERENCES foods,
  image_id integer NOT NULL REFERENCES images,

  created_at timestamp DEFAULT current_timestamp,
  updated_at timestamp DEFAULT current_timestamp,
  deleted_at timestamp DEFAULT null
);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE images, food_categories, foods, foods_images;
