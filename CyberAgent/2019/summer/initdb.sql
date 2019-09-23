create role cyberagent login password 'Kubernetes';

create database webapp with owner cyberagent;
alter database webapp set timezone to 'Asia/Tokyo';

\connect webapp
create table users (
       id         serial primary key,
       name       varchar(80) not null,
       email      text not null,
       created_at timestamp with time zone not null default current_timestamp,
       updated_at timestamp with time zone not null default current_timestamp
);

grant all privileges on users to cyberagent;
grant all privileges on all sequences in schema public to cyberagent;
