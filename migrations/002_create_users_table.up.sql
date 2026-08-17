Create Table users (
    id serial primary key,
    email varchar(150) not null unique,
    password_hash varchar(255) not null,
    role varchar(20) not null default 'USER',
    created_at TIMESTAMP not null default current_timestamp
);