-- +goose Up
ALTER TABLE user_information DROP COLUMN profile_complited;
ALTER TABLE user_information ADD COLUMN profile_complited BOOLEAN GENERATED ALWAYS AS (
    weight IS NOT NULL
    AND height IS NOT NULL
    AND age IS NOT NULL
    AND gender IS NOT NULL
) STORED;

-- +goose Down
ALTER TABLE user_information DROP COLUMN profile_complited;
ALTER TABLE user_information ADD COLUMN profile_complited BOOLEAN NOT NULL DEFAULT false;
