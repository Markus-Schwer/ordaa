CREATE TABLE IF NOT EXISTS users (
    uuid UUID DEFAULT gen_random_uuid(),
    name VARCHAR(40) NOT NULL,
    PRIMARY KEY (uuid)
);

