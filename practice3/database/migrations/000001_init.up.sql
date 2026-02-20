create table if not exists users (
  id serial primary key,
  name varchar(255) not null,
  email varchar(255) unique not null default '',
  password varchar(255) not null default '',
  deleted_at timestamp default null
);
insert into users (name, email, password) values ('John Doe', 'john@example.com', 'secret');
