-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
PRAGMA foreign_keys=off;

DROP TABLE if exists user;
CREATE TABLE IF NOT EXISTS user( 
   id integer primary key autoincrement,
   email text not null,
   name text not null,
   password text not null,
   created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);


DROP TABLE if exists conversation;

create table if not exists conversation (
    id integer primary key autoincrement,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

drop table if exists user_conversation;

create table if not exists user_conversation (
    user_id integer not null,
    conversation_id integer not null,

    constraint "user_conversation_user_fk_id" FOREIGN KEY (user_id) REFERENCES user (id) on delete restrict on update cascade,
    constraint "user_conversation_conversation_fk_id" FOREIGN KEY (conversation_id) REFERENCES conversation (id) on delete cascade on update cascade
);


drop table if exists messages;

create table if not exists messages (
    id integer not null,
    conversation_id integer not null,
    user_id integer not null,
    content text not null,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,


    constraint "message_user_fk_id" FOREIGN KEY (user_id) REFERENCES user (id) on delete restrict on update cascade,
    constraint "message_conversation_fk_id" FOREIGN KEY (conversation_id) REFERENCES conversation (id) on delete cascade on update cascade
);


PRAGMA foreign_key_check;
PRAGMA foreign_keys=ON;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
