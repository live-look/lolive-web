-- Deploy camforchat:convert_viewers_broadcast_id_to_uuid to pg

BEGIN;

  DROP INDEX index_viewers_broadcast_id_user_id;
  ALTER TABLE viewers ALTER COLUMN broadcast_id TYPE varchar(36);
  -- One viewer can view only one broadcast
  CREATE UNIQUE INDEX index_viewers_broadcast_id_user_id ON viewers (broadcast_id, user_id);

COMMIT;
