-- +goose Up
-- +goose StatementBegin


-- this will be the new setup here
-- disable foreign key constraint check
PRAGMA foreign_keys=off;


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
    constraint "todo_user_fk_id" FOREIGN KEY (user_id) REFERENCES user (id) on delete restrict on update cascade
);

PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
