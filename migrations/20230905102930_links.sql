-- +goose Up
-- +goose StatementBegin


-- this will be the new setup here
-- disable foreign key constraint check
PRAGMA foreign_keys=off;

-- start a transaction
-- BEGIN TRANSACTION;

DROP TABLE if exists todo;

-- setup both our tables
CREATE TABLE IF NOT EXISTS user( 
   id integer primary key autoincrement,
   first_name text,
   last_name text
);

-- create todo
create table if not exists todo (
    id integer primary key autoincrement,
    description text not null,
    user_id integer,
    FOREIGN KEY (user_id) 
      REFERENCES user (id)
);

PRAGMA foreign_keys=on;


-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
