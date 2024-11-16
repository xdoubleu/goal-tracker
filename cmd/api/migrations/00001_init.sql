-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS goals (
    id uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid NOT NULL,
    name varchar(255) NOT NULL,
    target_value integer,
    source_id integer NOT NULL,
    type_id integer NOT NULL,
    state varchar(255) NOT NULL
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
DROP TABLE IF EXISTS progress;
-- +goose StatementEnd
