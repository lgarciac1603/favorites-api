BEGIN;

CREATE TABLE user_favorites (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    crypto_id VARCHAR(50) NOT NULL,
    crypto_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    UNIQUE(user_id, crypto_id)
);

CREATE INDEX idx_user_favorites_user_id ON user_favorites(user_id);

COMMIT;