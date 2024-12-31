-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS goals (
    id varchar(255) NOT NULL PRIMARY KEY,
    parent_id varchar(255),
    name varchar(255) NOT NULL,
    is_linked boolean NOT NULL, 
    target_value integer,
    type_id integer,
    state varchar(255) NOT NULL,
    due_time timestamp,
    "order" integer NOT NULL
);

CREATE TABLE IF NOT EXISTS progress (
    type_id integer NOT NULL,
    date timestamp NOT NULL,
    value integer NOT NULL,
    PRIMARY KEY(type_id, date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS progress;
-- +goose StatementEnd
