-- Deploy camforchat:convert_broadcast_id_to_uuid to pg

BEGIN;

-- XXX Add DDLs here.
ALTER TABLE broadcasts DROP CONSTRAINT fk_broadcasts_user_id;
ALTER TABLE viewers DROP CONSTRAINT fk_viewers_broadcast_id;

ALTER TABLE broadcasts DROP COLUMN id;
DROP SEQUENCE IF EXISTS broadcasts_id_seq;

ALTER TABLE broadcasts ADD COLUMN id varchar(36);
CREATE UNIQUE INDEX pk_broadcasts_id ON broadcasts(id);
ALTER TABLE broadcasts ADD CONSTRAINT pk_broadcasts_id PRIMARY KEY USING INDEX pk_broadcasts_id;

COMMIT;
