// @generated automatically by Diesel CLI.

diesel::table! {
    menu_items (id) {
        id -> Integer,
        short_name -> Text,
        name -> Text,
        menu_id -> Integer,
        price -> Integer,
    }
}

diesel::table! {
    menus (id) {
        id -> Integer,
        name -> Text,
        url -> Nullable<Text>,
    }
}

diesel::table! {
    order_items (id) {
        id -> Integer,
        menu_item_id -> Integer,
        user -> Integer,
        paid -> Bool,
        order_id -> Integer,
        price -> Integer,
    }
}

diesel::table! {
    orders (id) {
        id -> Integer,
        order_deadline -> Nullable<Integer>,
        eta -> Nullable<Integer>,
        initiator -> Integer,
        sugar_person -> Nullable<Integer>,
        state -> Text,
        menu_id -> Integer,
    }
}

diesel::table! {
    users (id) {
        id -> Integer,
        name -> Text,
    }
}

diesel::joinable!(menu_items -> menus (menu_id));
diesel::joinable!(order_items -> menu_items (menu_item_id));
diesel::joinable!(order_items -> orders (order_id));
diesel::joinable!(order_items -> users (user));
diesel::joinable!(orders -> users (initiator));
diesel::joinable!(orders -> menus (menu_id));

diesel::allow_tables_to_appear_in_same_query!(
    menu_items,
    menus,
    order_items,
    orders,
    users,
);
