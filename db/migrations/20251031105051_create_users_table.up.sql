create table if not exists users (
    id integer primary key autoincrement,
    email text unique not null,
    google_id text unique,
    created_at datetime default current_timestamp,
    updated_at datetime default current_timestamp,
    last_login_at datetime
);

