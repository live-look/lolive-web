-- Deploy camforchat:create_broadcasts to pg

CREATE TYPE broadcast_state AS ENUM (
  'online',
  'offline',
  'private'
);

CREATE TYPE viewer_state AS ENUM (
  'joined',
  'exited'
);

BEGIN;

  CREATE TABLE broadcasts (
    id bigserial NOT NULL PRIMARY KEY,
    user_id bigint NOT NULL,
    created_at timestamp with time zone NOT NULL,
    state broadcast_state NOT NULL DEFAULT 'offline'
  );

  CREATE INDEX index_broadcasts_user_id ON broadcasts (user_id);

  CREATE TABLE viewers (
    id bigserial NOT NULL PRIMARY KEY,
    broadcast_id bigint NOT NULL,
    user_id bigint NOT NULL,
    state viewer_state NOT NULL DEFAULT 'joined',
    joined_at timestamp with time zone NOT NULL,
    exited_at timestamp with time zone
  );

  CREATE UNIQUE INDEX index_viewers_broadcast_id_user_id ON viewers (broadcast_id, user_id);
  CREATE INDEX index_viewers_user_id ON viewers (user_id);

  ALTER TABLE broadcasts ADD CONSTRAINT fk_broadcasts_user_id FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;

  ALTER TABLE viewers ADD CONSTRAINT fk_viewers_broadcast_id FOREIGN KEY (broadcast_id) REFERENCES broadcasts(id)
    ON DELETE CASCADE;

  ALTER TABLE viewers ADD CONSTRAINT fk_viewers_user_id FOREIGN KEY (user_id) REFERENCES users(id)
    ON DELETE CASCADE;
COMMIT;
