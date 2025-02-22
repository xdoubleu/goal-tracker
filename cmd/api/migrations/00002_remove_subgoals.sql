-- +goose Up
-- +goose StatementBegin
DELETE FROM goals WHERE parent_id IS NOT null;
ALTER TABLE goals DROP COLUMN parent_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
