CREATE TABLE event_views (
    id         BIGSERIAL PRIMARY KEY,
    event_id   BIGINT NOT NULL REFERENCES romantic_events(id) ON DELETE CASCADE,
    visitor_id TEXT   NOT NULL,
    device     TEXT   NOT NULL DEFAULT '',
    os         TEXT   NOT NULL DEFAULT '',
    browser    TEXT   NOT NULL DEFAULT '',
    ip         TEXT   NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX idx_event_views_event_id      ON event_views(event_id);
CREATE INDEX idx_event_views_event_visitor ON event_views(event_id, visitor_id);
