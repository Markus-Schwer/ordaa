pub mod menu;
pub mod order;

use actix_web::web;

pub fn configure(cfg: &mut web::ServiceConfig) {
    menu::services_menu(cfg);
    order::services_order(cfg);
}
