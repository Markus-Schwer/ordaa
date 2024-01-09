CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(40)
);

CREATE TABLE IF NOT EXISTS order_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    menu VARCHAR(255) NOT NULL,
    item_id VARCHAR(255) NOT NULL,
    user INTEGER NOT NULL,
    paid BOOLEAN NOT NULL,
    order_id INTEGER NOT NULL,
    FOREIGN KEY(menu, item_id) REFERENCES menu_items(menu, id),
    FOREIGN KEY(user) REFERENCES users(id),
    FOREIGN KEY(order_id) REFERENCES orders(id)
);

CREATE TABLE IF NOT EXISTS orders (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    order_deadline INTEGER,
    eta INTEGER,
    initiator INTEGER NOT NULL,
    sugar_person INTEGER,
    state VARCHAR(40) NOT NULL,
    FOREIGN KEY(initiator) REFERENCES users(id),
    FOREIGN KEY(sugar_person) REFERENCES users(id)
);
