-- Revert camforchat:messages from pg

BEGIN;

  DROP TABLE broadcast_messages;

COMMIT;
