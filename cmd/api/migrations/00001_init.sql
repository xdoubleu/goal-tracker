-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS goals (
    id varchar(255) NOT NULL PRIMARY KEY,
    parent_id varchar(255),
    user_id uuid NOT NULL,
    name varchar(255) NOT NULL,
    is_linked boolean NOT NULL, 
    target_value integer,
    type_id integer,
    state varchar(255) NOT NULL,
    due_time timestamp
);

CREATE TABLE IF NOT EXISTS progress (
    id serial4 PRIMARY KEY,
    goal_id varchar(255) NOT NULL REFERENCES goals ON DELETE CASCADE,
    value integer NOT NULL,
    created_at timestamp NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS progress;
-- +goose StatementEnd
