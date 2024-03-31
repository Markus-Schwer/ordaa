CREATE TABLE IF NOT EXISTS password_users (
    uuid UUID DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    username TEXT NOT NULL,
    password TEXT NOT NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_password_users_user FOREIGN KEY(user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS matrix_users (
    uuid UUID DEFAULT gen_random_uuid(),
    user_uuid UUID NOT NULL,
    username TEXT NOT NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_matrix_users_user FOREIGN KEY(user_uuid) REFERENCES users(uuid) ON DELETE CASCADE
);
