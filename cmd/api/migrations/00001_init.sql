-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS goals (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    name varchar(255) NOT NULL,
    description text,
    date timestamp,
    value integer,
    source_id integer NOT NULL,
    type_id integer NOT NULL,
    score integer NOT NULL CHECK (score > 0),
    state_id integer NOT NULL
);

CREATE TABLE IF NOT EXISTS progress (
    id serial4 PRIMARY KEY,
    goal_id uuid NOT NULL REFERENCES goals ON DELETE CASCADE,
    value integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goals;
-- +goose StatementEnd
