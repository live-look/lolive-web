-- Deploy camforchat:accounts to pg

BEGIN;

  CREATE TABLE accounts (
    id bigserial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    total numeric(12, 2) NOT NULL DEFAULT 0
  );

  CREATE UNIQUE INDEX index_accounts_user_id ON accounts (user_id);

COMMIT;
