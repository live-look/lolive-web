-- Deploy camforchat:fk_viewers_broadcast_id to pg

BEGIN;

  ALTER TABLE viewers ADD CONSTRAINT fk_viewers_broadcast_id FOREIGN KEY (broadcast_id) REFERENCES broadcasts(id)
    DEFERRABLE INITIALLY DEFERRED;

COMMIT;
