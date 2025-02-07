create table if not exists apple_tokens
(
    id            int auto_increment
        primary key,
    user_id       char(255)                           not null,
    created_at    timestamp default CURRENT_TIMESTAMP not null,
    access_token  text                                not null,
    token_type    char(128)                           null,
    expires_in    int       default 3600              not null,
    refresh_token text                                not null,
    id_token      text                                not null,
    constraint apple_tokens_pk
        unique (user_id)
);

create table users
(
    id         int auto_increment
        primary key,
    user_id    char(255)                           not null,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    auth_type  int       default 0                 not null,
    constraint users_pk
        unique (user_id)
);
