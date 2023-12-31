use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct MenuItem {
    id: String,
    name: String,
    price: usize,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Menu {
    pub items: Vec<MenuItem>,
}

pub mod filters {
    use sqlx::SqlitePool;
    use warp::{http::StatusCode, Filter};

    use crate::filters::{json_body, with_db};

    use super::Menu;

    pub fn all(
        pool: &SqlitePool,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        update_menu(pool)
        // todos_list(db.clone())
        //     .or(todos_create(db.clone()))
        //     .or(todos_update(db.clone()))
        //     .or(todos_delete(db))
    }

    fn update_menu(
        pool: &SqlitePool,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu" / String)
            .and(warp::put())
            .and(json_body())
            .and(with_db(pool.clone()))
            .and_then(|name: String, menu: Menu, pool: SqlitePool| async move {
                let mut conn = pool.acquire().await.unwrap();
                for it in menu.items {
                    sqlx::query("INSERT INTO MENU_ITEM VALUES (?1, ?2, ?3, ?4)")
                        .bind(name.clone())
                        .bind(it.id)
                        .bind(it.name)
                        .bind(it.price as i64)
                        .execute(&mut *conn)
                        .await
                        .unwrap();
                    // statement.bind_iter::<_, (_, Value)>([
                    //     (":menu", name.clone().into()),
                    //     (":id", it.id.to_owned().into()),
                    //     (":name", it.name.to_owned().into()),
                    //     (":price", ( it.price.to_owned() as i64 ).into()),
                    // ]).unwrap();
                }
                Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
                // match state.menus.write() {
                //     Ok(mut locked_menus) => {
                //
                //             Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
                //         // if locked_menus.contains_key(&name) {
                //         //     locked_menus.insert(name, menu);
                //         //     Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
                //         // } else {
                //         //     locked_menus.insert(name, menu);
                //         //     Ok(StatusCode::CREATED)
                //         // }
                //     }
                //     Err(_) => Ok(StatusCode::SERVICE_UNAVAILABLE),
                // }
            })
    }
}
