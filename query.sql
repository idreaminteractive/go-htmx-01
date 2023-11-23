

-- name: CreateUser :one
insert into user (
  password, email, handle
) values (? , ?, ?) returning *;


-- name: GetUserByEmail :one
select * from user 
where email = ? limit 1; 


-- name: GetAllUsers :many
select * from user;

-- name: CreateConversation :one
insert into conversation (topic) values ("") returning  *;

-- name: LinkUserToConversation :one
insert into user_conversation (user_id, conversation_id) values (?, ?) returning *;

-- name: GetConversationsForUser :many 
select m.content, m.id, m.user_id, u.handle, u.id 
from user_conversation uc, messages m, user u 
where uc.user_id = ? and uc.conversation_id = m.conversation_id and uc.user_id = u.id;

-- name: GetConversationsList :many
select
  uc.conversation_id, 
u.handle,
u.id as user_id,
  json_group_array(json_object(
    'message_id', m.id,
    'content', m.content,
    'user_id', m.user_id,
    'handle', m.handle,
    'created_at', m.created_at
   )) as conversation_messages
   
from
  user_conversation uc
    join (select messages.id, messages.created_at, messages.conversation_id, messages.content, messages.user_id, u.handle from messages, user u where u.id = messages.user_id  order by messages.created_at desc) as m on m.conversation_id = uc.conversation_id
    -- get the other user in the conversation who is NOT me.
    
    join user u on uc.user_id = u.id 
    where uc.user_id = ?
    group by uc.conversation_id
order by
  uc.conversation_id
limit
  10;


-- name: CreateMessage :one
insert into messages (user_id, conversation_id, content) values (?, ?, ?) returning *;


-- name: GetOtherConversationUser :one
select u.id, u.handle from user u, user_conversation uc where u.id = uc.user_id and uc.conversation_id=? and u.id != ? limit 1;