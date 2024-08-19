create type user_role as enum('client', 'moderator');
create type apartment_status as enum('created', 'approved', 'decline', 'on moderation');

--- Создание таблицы пользователей
create table users (
    user_id uuid primary key,
    mail varchar(255) not null unique,
    password text not null,
    role user_role not null
);

--- Создание таблицы домов
create table houses (
    house_id serial primary key,
    address text not null,
    construction_date int,
    developer text,
    created_at timestamp without time zone default current_timestamp,
    updated_at timestamp without time zone
);

--- Создание таблицы квартир
create table apartments (
    apartment_id serial primary key,
    house_id int references houses(house_id) on delete cascade,
    moderator_id uuid references users(user_id),
    user_id uuid references users(user_id),
    price int not null,
    rooms int not null,
    status apartment_status not nullx
);
