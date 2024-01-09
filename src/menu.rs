use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow, sqlx::Type, Clone)]
pub struct MenuItem {
    pub id: String,
    pub name: String,
    pub price: i64,
    pub menu: String
}

#[derive(Debug, Serialize, Deserialize, sqlx::FromRow, Clone)]
pub struct Menu {
    pub name: String,
    #[sqlx(skip)]
    pub items: Vec<MenuItem>,
}

pub mod filters {
    use serde::{Deserialize, Serialize};
    use warp::{http::StatusCode, Filter};

    use crate::{
        filters::{json_body, with_db, with_searcher_ctx},
        search::SearchContextReader, db::Db,
    };

    use super::Menu;

    pub fn all(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        update_menu(db.clone(), ctx.clone()).or(get_menu(db.clone(), ctx.clone()))
    }

    #[derive(Deserialize, Serialize, Debug)]
    struct FuzzyParam {
        search_string: String,
    }

    fn get_menu(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        let opt_query = warp::query::<FuzzyParam>()
            .map(Some)
            .or_else(|_| async { Ok::<(Option<FuzzyParam>,), std::convert::Infallible>((None,)) });
        warp::path!("menu" / String)
            .and(warp::get())
            .and(opt_query)
            .and(with_db(db.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(
                |name: String,
                 query: Option<FuzzyParam>,
                 db: Db,
                 ctx: SearchContextReader| async move {
                    let items = if let Some(param) = query {
                        let ids = ctx.fuzz_menu_item_ids(param.search_string.as_str());
                        db.get_items_by_id(ids, name).await
                    } else {
                        db.all_items(name).await
                    };
                    Ok::<warp::reply::Json, warp::Rejection>(warp::reply::json(&items))
                },
            )
    }

    fn update_menu(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("menu" / String)
            .and(warp::put())
            .and(json_body())
            .and(with_db(db.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(
                |_: String, menu: Menu, db: Db, ctx: SearchContextReader| async move {
                    ctx.index_write_sender.send(menu.clone()).unwrap();
                    db.insert_menu(menu).await;
                    Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
                },
            )
    }
}
