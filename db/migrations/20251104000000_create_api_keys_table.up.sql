create table if not exists api_keys (
    id text primary key,
    user_id integer not null,
    name text not null,
    key_hash text not null,
    key_prefix text not null,
    key_suffix text not null,
    created_at datetime not null default current_timestamp,
    last_used_at datetime,
    foreign key (user_id) references users(id) on delete cascade
);

create index idx_api_keys_user_id on api_keys(user_id);
create index idx_api_keys_key_hash on api_keys(key_hash);
