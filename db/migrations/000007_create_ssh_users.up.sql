CREATE TABLE IF NOT EXISTS ssh_users (
    uuid UUID DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    public_key TEXT NOT NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_matrix_users_user FOREIGN KEY(user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);
INSERT INTO ssh_users (user_uuid, public_key) SELECT users.uuid, users.public_key FROM users WHERE users.public_key IS NOT NULL AND users.uuid NOT IN (SELECT ssh_users.user_uuid FROM ssh_users);
ALTER TABLE users DROP COLUMN IF EXISTS public_key;
