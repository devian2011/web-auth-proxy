CREATE TABLE roles
(
    id     varchar(255) primary key,
    name   varchar(255) not null,
    code   varchar(255) not null
);

CREATE TABLE users
(
    id            varchar(1024) primary key,
    username      varchar(1024)  not null unique,
    password_hash varchar(2048) not null,
    email         varchar(512),
    phone         varchar(512)
);

CREATE TABLE roles_for_user
(
    user_id varchar(255) not null references users (id),
    role_id varchar(255) not null references roles (id)
);

CREATE TABLE sessions
(
    id         varchar(255) primary key,
    user_id    varchar(255) references users (id),
    start_in   bigint default 0,
    expires_in bigint default 0,
    token      varchar(255) not null
);

INSERT INTO roles (id, name, code, source)
VALUES ('1', 'Admin', 'ADMIN', 'proxy'),
       ('2', 'Manager', 'MANAGER', 'proxy'),
       ('3', 'User', 'USER', 'proxy');