BEGIN;

CREATE TABLE IF NOT EXISTS posts (
    id BIGSERIAL PRIMARY KEY,
    newsletter_id UUID,
    heading VARCHAR(32) NOT NULL,
    body VARCHAR(255) NOT NULL,
    FOREIGN KEY (newsletter_id) REFERENCES newsletters (id)
);


COMMIT;
