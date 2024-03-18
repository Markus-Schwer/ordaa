CREATE TABLE users (
    uuid UUID DEFAULT gen_random_uuid(),
    name VARCHAR(40) NOT NULL,
    PRIMARY KEY (uuid)
);

CREATE TABLE orders (
    uuid UUID DEFAULT gen_random_uuid(),
    order_deadline INTEGER,
    eta INTEGER,
    initiator UUID NOT NULL,
    sugar_person UUID,
    state VARCHAR(40) NOT NULL,
    menu_uuid UUID NOT NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_orders_initiator FOREIGN KEY(initiator) REFERENCES users(uuid),
    CONSTRAINT fk_orders_sugar_person FOREIGN KEY(sugar_person) REFERENCES users(uuid),
    CONSTRAINT fk_orders_menu FOREIGN KEY(menu_uuid) REFERENCES menus(uuid)
);

CREATE TABLE order_items (
    uuid UUID DEFAULT gen_random_uuid(),
    menu_item_uuid UUID NOT NULL,
    order_user UUID NOT NULL,
    paid BOOLEAN NOT NULL DEFAULT FALSE,
    order_uuid UUID NOT NULL,
    price INTEGER NOT NULL,
    PRIMARY KEY (uuid),
    CONSTRAINT fk_order_items_menu_item FOREIGN KEY(menu_item_uuid) REFERENCES menu_items(uuid),
    CONSTRAINT fk_order_items_user FOREIGN KEY(order_user) REFERENCES users(uuid),
    CONSTRAINT fk_order_items_order FOREIGN KEY(order_uuid) REFERENCES orders(uuid)
);
