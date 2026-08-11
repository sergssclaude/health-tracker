-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE, 
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(100),
    weight INT,
    height INT,
    age INT,
    gender INT,
    goal INT,
    daily_calorie_target INT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP
);


-- +goose Down
DROP TABLE users;
