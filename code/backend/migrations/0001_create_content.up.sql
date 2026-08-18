CREATE TABLE contents (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    value      TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT ck_contents_value_not_blank CHECK (length(value) > 0)
);

CREATE UNIQUE INDEX uq_contents_singleton ON contents ((true));

INSERT INTO contents (value) VALUES ('Hello Word');
