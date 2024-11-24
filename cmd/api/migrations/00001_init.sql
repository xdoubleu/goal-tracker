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

CREATE TABLE IF NOT EXISTS steam_games (
    app_id integer PRIMARY KEY,
	name varchar(255),
	playtime_2weeks integer,
	playtime_forever integer,
	img_icon_url varchar(255), 
	img_logo_url varchar(255),
	has_community_visible_stats boolean   
);

CREATE TABLE IF NOT EXISTS steam_achievements (
    id serial4 PRIMARY KEY,
    app_id integer NOT NULL REFERENCES steam_games ON DELETE CASCADE,
    apiname varchar(255),
	name varchar(255),
	description varchar(255)  
);

CREATE TABLE IF NOT EXISTS steam_achievements_user (
    steam_achievement_id serial4 NOT NULL REFERENCES steam_achievements ON DELETE CASCADE,
    user_id uuid NOT NULL,
    achieved boolean,
	unlocktime timestamp,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goals;
DROP TABLE IF EXISTS progress;
DROP TABLE IF EXISTS steam_games;
DROP TABLE IF EXISTS steam_achievements;
DROP TABLE IF EXISTS steam_achievements_user;
-- +goose StatementEnd
