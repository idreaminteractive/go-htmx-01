-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
PRAGMA foreign_keys=off;


-- setup both our tables
alter TABLE user add column password text;


PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
