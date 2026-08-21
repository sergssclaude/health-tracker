-- +goose Up
CREATE TABLE food_items (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    calories_per_100g INT NOT NULL,
    protein_per_100g  INT,
    fats_per_100g  INT,
    carbs_per_100g  INT
);

CREATE TABLE food_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    food_item_id INT NOT NULL REFERENCES food_items(id),
    amount_grams INT NOT NULL, 
    meal_type VARCHAR(20) NOT NULL CHECK (meal_type IN('breakfast', 'lunch', 'dinner', 'snack')),
    logged_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE food_logs;
DROP TABLE food_items;

