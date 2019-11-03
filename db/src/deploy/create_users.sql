-- Deploy camforchat:create_users to pg

BEGIN;

  CREATE TABLE users (
    id bigserial NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    password varchar(1024) NOT NULL,
    confirmed boolean NOT NULL DEFAULT 'f',
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
  );

  CREATE UNIQUE INDEX uniq_users_email ON users (lower(email));
COMMIT;
