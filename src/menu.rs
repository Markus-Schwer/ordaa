use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow, Clone)]
pub struct MenuItem {
    pub id: String,
    pub name: String,
    pub price: i64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Menu {
    pub name: String,
    pub items: Vec<MenuItem>,
}

pub mod filters {
    use serde::{Deserialize, Serialize};
    use sqlx::SqlitePool;
    use warp::{http::StatusCode, Filter};

    use crate::{
        filters::{json_body, with_db, with_searcher_ctx},
        search::SearchContextReader,
    };

    use super::{Menu, MenuItem};

    pub fn all(
        pool: &SqlitePool,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        update_menu(pool, ctx.clone()).or(get_menu(pool, ctx.clone()))
    }

    #[derive(Deserialize, Serialize, Debug)]
    struct FuzzyParam {
        search_string: String,
    }

    fn get_menu(
        pool: &SqlitePool,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        let opt_query = warp::query::<FuzzyParam>()
            .map(Some)
            .or_else(|_| async { Ok::<(Option<FuzzyParam>,), std::convert::Infallible>((None,)) });
        warp::path!("menu" / String)
            .and(warp::get())
            .and(opt_query)
            .and(with_db(pool.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(
                |name: String,
                 query: Option<FuzzyParam>,
                 pool: SqlitePool,
                 ctx: SearchContextReader| async move {
                    let mut conn = pool.acquire().await.unwrap();
                    let items = if let Some(param) = query {
                        let ids = ctx.fuzz_menu_item_ids(param.search_string.as_str());
                        let mut matches: Vec<MenuItem> = Vec::new();
                        for id in ids {
                            matches.push(
                                sqlx::query_as::<_, MenuItem>(
                                    "SELECT id, name, price FROM MENU_ITEM WHERE id = ?1",
                                )
                                .bind(id)
                                .fetch_one(&mut *conn)
                                .await
                                .unwrap(),
                            );
                        }
                        matches
                    } else {
                        sqlx::query_as::<_, MenuItem>(
                            "SELECT id, name, price FROM MENU_ITEM WHERE menu = ?1",
                        )
                        .bind(name.clone())
                        .fetch_all(&mut *conn)
                        .await
                        .unwrap()
                    };
                    Ok::<warp::reply::Json, warp::Rejection>(warp::reply::json(&items))
                },
            )
    }

    fn update_menu(
        pool: &SqlitePool,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu" / String)
            .and(warp::put())
            .and(json_body())
            .and(with_db(pool.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(
                |name: String, menu: Menu, pool: SqlitePool, ctx: SearchContextReader| async move {
                    ctx.index_write_sender.send(menu.clone()).unwrap();
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
                    }
                    Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
                },
            )
    }
}
