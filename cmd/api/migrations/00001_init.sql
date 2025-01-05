-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp"; -- noqa: L057

CREATE TABLE IF NOT EXISTS states (
    id varchar(255) NOT NULL,
    user_id varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    "order" integer NOT NULL,
    PRIMARY KEY (id, user_id)
);

CREATE TABLE IF NOT EXISTS goals (
    id varchar(255) NOT NULL,
    user_id varchar(255) NOT NULL,
    parent_id varchar(255),
    name varchar(255) NOT NULL,
    state_id varchar(255) NOT NULL,
    "order" integer NOT NULL,
    is_linked boolean NOT NULL DEFAULT false,
    source_id integer,
    type_id integer,
    target_value integer,
    period integer,
    due_time timestamp,
    config json,
    PRIMARY KEY (id, user_id)
);

CREATE TABLE IF NOT EXISTS progress (
    type_id integer NOT NULL,
    user_id varchar(255) NOT NULL,
    date timestamp NOT NULL,
    value varchar(255) NOT NULL,
    PRIMARY KEY (type_id, user_id, date)
);

CREATE TABLE IF NOT EXISTS list_items (
    id integer NOT NULL,
    user_id varchar(255) NOT NULL,
    goal_id varchar(255) NOT NULL,
    value varchar(255) NOT NULL,
    completed boolean NOT NULL,
    PRIMARY KEY (id, user_id, goal_id)
);

CREATE TABLE IF NOT EXISTS goodreads_books (
    id integer NOT NULL,
    user_id varchar(255) NOT NULL,
    shelf varchar(255) NOT NULL,
    tags varchar(255) [] NOT NULL,
    title varchar(255) NOT NULL,
    author varchar(255) NOT NULL,
    dates_read timestamp [],
    PRIMARY KEY (id, user_id)
);

CREATE TABLE IF NOT EXISTS steam_games (
    id integer NOT NULL,
    user_id varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    is_delisted boolean NOT NULL DEFAULT false,
    PRIMARY KEY (id, user_id)
);

CREATE TABLE IF NOT EXISTS steam_achievements (
    name varchar(255) NOT NULL,
    user_id varchar(255) NOT NULL,
    game_id integer NOT NULL,
    achieved boolean NOT NULL DEFAULT false,
    unlock_time timestamp,
    PRIMARY KEY (name, user_id, game_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS states;
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS progress;
DROP TABLE IF EXISTS list_items;
DROP TABLE IF EXISTS goodreads_books;
DROP TABLE IF EXISTS steam_games;
DROP TABLE IF EXISTS steam_achievements;
-- +goose StatementEnd
