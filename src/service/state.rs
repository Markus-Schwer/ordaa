use crate::entity::{db::Db, search::SearchContextReader};

#[derive(Clone)]
pub struct AppState {
    pub db: Db,
    pub search: SearchContextReader
}
