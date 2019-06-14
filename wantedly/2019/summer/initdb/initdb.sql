create role wantedly login;

create database webapp with owner wantedly;
alter database webapp set timezone to 'Asia/Tokyo';

\connect webapp
create table users (
       id         serial primary key,
       name       varchar(80) not null,
       email      text not null,
       created_at timestamp with time zone not null default current_timestamp,
       updated_at timestamp with time zone not null default current_timestamp
);

grant all privileges on users to wantedly;
grant all privileges on all sequences in schema public to wantedly;
