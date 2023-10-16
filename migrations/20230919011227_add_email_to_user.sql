-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
PRAGMA foreign_keys=off;


-- setup both our tables
alter TABLE user add column email text NOT NULL ;
create unique index  if not exists user_email_idx on user (email);


PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
