use crate::dto::MenuItemDto;

pub struct Menu {
    pub id: i32,
    pub name: String,
    pub url: Option<String>,
    pub items: Vec<MenuItemDto>
}

pub struct MenuItem {
    pub id: i32,
    pub short_name: String,
    pub name: String,
    pub price: i32,
}

pub mod filters {
    use serde::{Deserialize, Serialize};
    use warp::{http::StatusCode, Filter};

    use crate::{
        filters::{json_body, with_db, with_searcher_ctx},
        search::SearchContextReader, db::Db, dto::NewMenuDto,
    };

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
        warp::path!("api" / "menu" / i32)
            .and(warp::get())
            .and(opt_query)
            .and(with_db(db.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(
                |id: i32,
                 query: Option<FuzzyParam>,
                 db: Db,
                 ctx: SearchContextReader| async move {
                    let items = if let Some(param) = query {
                        let ids = ctx.fuzz_menu_item_ids(param.search_string.as_str());
                        db.get_items_by_id(ids, id)
                    } else {
                        db.all_items(id)
                    };
                    Ok::<warp::reply::Json, warp::Rejection>(warp::reply::json(&items))
                },
            )
    }

    fn update_menu(
        db: Db,
        ctx: SearchContextReader,
    ) -> impl Filter<Extract = (impl warp::Reply,), Error = warp::Rejection> + Clone {
        warp::path!("api" / "menu" / String)
            .and(warp::put())
            .and(json_body())
            .and(with_db(db.clone()))
            .and(with_searcher_ctx(ctx))
            .and_then(
                |_: String, new_menu: NewMenuDto, db: Db, ctx: SearchContextReader| async move {
                    let menu = db.insert_menu(new_menu);
                    ctx.index_write_sender.send(menu).unwrap();
                    Ok::<StatusCode, warp::Rejection>(StatusCode::OK)
                },
            )
    }
}
