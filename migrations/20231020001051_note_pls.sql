-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
PRAGMA foreign_keys=off;

-- Drop todo
DROP TABLE if exists todo;

-- create notes
create table if not exists note (
    id integer primary key autoincrement,
    content text not null,
    user_id integer not null,
    is_public BOOLEAN  NOT NULL DEFAULT false,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    constraint "notes_user_fk_id" FOREIGN KEY (user_id) REFERENCES user (id) on delete restrict on update cascade
);


PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
