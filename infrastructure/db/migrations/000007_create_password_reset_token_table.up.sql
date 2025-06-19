create table password_reset_tokens
(
    token      text primary key,
    user_id    uuid                     NOT NULL,
    expires_at timestamp with time zone not null,

    constraint fk_user
        foreign key (user_id)
            references users (id)
            on delete cascade
);

create index idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);