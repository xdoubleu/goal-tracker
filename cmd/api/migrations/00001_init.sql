-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS states (
    id varchar(255) PRIMARY KEY,
    name varchar(255) NOT NULL,
    "order" integer NOT NULL
);

CREATE TABLE IF NOT EXISTS goals (
    id varchar(255) NOT NULL PRIMARY KEY,
    parent_id varchar(255),
    name varchar(255) NOT NULL,
    state_id varchar(255) NOT NULL REFERENCES states,
    "order" integer NOT NULL,
    is_linked boolean NOT NULL,
    type_id integer,
    target_value integer,
    due_time timestamp,
    config json
);

CREATE TABLE IF NOT EXISTS progress (
    type_id integer NOT NULL,
    date timestamp NOT NULL,
    value varchar(255) NOT NULL,
    PRIMARY KEY (type_id, date)
);

CREATE TABLE IF NOT EXISTS list_items (
    id integer NOT NULL,
    goal_id varchar(255) NOT NULL REFERENCES goals,
    value varchar(255) NOT NULL,
    completed boolean NOT NULL,
    PRIMARY KEY (id, goal_id)
);

CREATE TABLE IF NOT EXISTS goodreads_books (
    id integer PRIMARY KEY,
    shelf varchar(255) NOT NULL,
    tags varchar(255) [] NOT NULL,
    title varchar(255) NOT NULL,
    author varchar(255) NOT NULL,
    dates_read timestamp []
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS states;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS progress;
DROP TABLE IF EXISTS list_items;
DROP TABLE IF EXISTS goodreads_books;
-- +goose StatementEnd
