CREATE TABLE clicks (
  id         BIGSERIAL PRIMARY KEY,
  link_id    BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
  timestamp  TIMESTAMPTZ NOT NULL DEFAULT now(),
  referrer   TEXT
);
