-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
  id UUID NOT NULL PRIMARY KEY,
  name TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
