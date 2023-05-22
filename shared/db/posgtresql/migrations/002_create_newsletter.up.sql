BEGIN;

CREATE TABLE IF NOT EXISTS newsletters (
   id UUID PRIMARY KEY,
   user_id UUID,
   name VARCHAR(255) NOT NULL,
   FOREIGN KEY (user_id) REFERENCES users (id)
);


COMMIT;
