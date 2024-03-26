CREATE TABLE IF NOT EXISTS menus (
    uuid UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    url VARCHAR(255),
    PRIMARY KEY (uuid)
);

CREATE TABLE IF NOT EXISTS menu_items (
    uuid UUID DEFAULT gen_random_uuid(),
    short_name VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    menu_uuid UUID NOT NULL,
    price INTEGER NOT NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_menu_items_menu FOREIGN KEY(menu_uuid) REFERENCES menus(uuid) ON DELETE CASCADE
);

