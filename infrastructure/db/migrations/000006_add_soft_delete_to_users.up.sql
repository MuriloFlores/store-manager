alter table users
add column deleted_at timestamp with time zone null;

create index idx_users_deleted_at on users (deleted_at)