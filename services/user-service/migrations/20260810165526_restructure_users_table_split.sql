-- +goose Up
ALTER TABLE users DROP COLUMN weight;
ALTER TABLE users DROP COLUMN height;
ALTER TABLE users DROP COLUMN age;
ALTER TABLE users DROP COLUMN gender;
ALTER TABLE users DROP COLUMN goal;
ALTER TABLE users DROP COLUMN daily_calorie_target;

CREATE TABLE user_information (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    weight NUMERIC(5,2),
    height INT,
    age INT, 
    gender VARCHAR(10) CHECK (gender IN ('male', 'female')),
    daily_calorie_norm INT,
    profile_complited BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);

CREATE TABLE user_goal (
    user_id INT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    target_weight NUMERIC(5,2),
    calorie_goal INT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE user_goal;
DROP TABLE user_information;

ALTER TABLE users ADD COLUMN weight INT;
ALTER TABLE users ADD COLUMN height INT;
ALTER TABLE users ADD COLUMN age INT;
ALTER TABLE users ADD COLUMN gender INT;
ALTER TABLE users ADD COLUMN goal INT;
ALTER TABLE users ADD COLUMN daily_calorie_target INT;



