-- Deploy camforchat:add_confirm_fields to pg

BEGIN;

  ALTER TABLE users ADD COLUMN confirm_selector varchar(2000),
                    ADD COLUMN confirm_verifier varchar(2000);

COMMIT;
