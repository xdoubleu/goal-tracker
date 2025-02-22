-- +goose Up
-- +goose StatementBegin
DROP TABLE list_items;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS list_items (
    id integer NOT NULL,
    user_id varchar(255) NOT NULL,
    goal_id varchar(255) NOT NULL,
    value varchar(255) NOT NULL,
    completed boolean NOT NULL,
    PRIMARY KEY (id, user_id, goal_id)
);
-- +goose StatementEnd
