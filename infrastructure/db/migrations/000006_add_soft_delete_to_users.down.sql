drop index if exists idx_users_deleted_at;
alter table users drop column deleted_at;