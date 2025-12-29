-- +goose Up
-- +goose StatementBegin
CREATE TYPE preferences_enum AS ENUM ('CARDIO', 'WEIGHT');
CREATE TYPE weight_unit_enum AS ENUM ('KG', 'LBS');
CREATE TYPE height_unit_enum AS ENUM ('CM', 'INCH');

CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  name VARCHAR(60),
  preference preferences_enum,
  weight_unit weight_unit_enum,
  height_unit height_unit_enum,
  weight SMALLINT,
  height SMALLINT,
  image_uri TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  is_active bit default B'1'
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
DROP TYPE IF EXISTS preferences_enum;
DROP TYPE IF EXISTS weight_unit_enum;
DROP TYPE IF EXISTS height_unit_enum;
-- +goose StatementEnd
