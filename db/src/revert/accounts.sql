-- Revert camforchat:accounts from pg

BEGIN;

  DROP TABLE accounts;

COMMIT;
