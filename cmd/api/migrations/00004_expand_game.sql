-- +goose Up
-- +goose StatementBegin
ALTER TABLE steam_games ADD COLUMN completion_rate VARCHAR(255) DEFAULT '';
ALTER TABLE steam_games ADD COLUMN contribution VARCHAR(255) DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE steam_games DROP COLUMN completion_rate;
ALTER TABLE steam_games DROP COLUMN contribution;
-- +goose StatementEnd
