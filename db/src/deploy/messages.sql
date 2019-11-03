-- Deploy camforchat:messages to pg

BEGIN;

  CREATE TABLE broadcast_messages (
    id bigserial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    broadcast_id bigint NOT NULL,
    is_system boolean NOT NULL DEFAULT 'f',
    created_at timestamp with time zone NOT NULL,
    content text
  );

COMMIT;
