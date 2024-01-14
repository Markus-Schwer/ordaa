mod boundary;
mod service;
mod entity;
mod frontend;

use entity::db::Db;
use entity::search::init_search_index;
use actix_web::{web, App, HttpServer};
use service::state::AppState;

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let (search_writer, search_reader) = init_search_index();
    let index_writer_handle = search_writer.start_index_writer(search_reader.clone());

    let db = Db::new();
    db.init_schema();
    let actix_handle = HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(AppState {
                db: db.clone(),
                search: search_reader.clone()
            }))
            .service(web::scope("/api").configure(boundary::menu::services_menu))
            .configure(frontend::services_frontend)
    })
    .bind(("127.0.0.1", 8080))
    .unwrap().run();

    let (_, _) = tokio::join!(index_writer_handle, actix_handle);
    Ok(())
}
