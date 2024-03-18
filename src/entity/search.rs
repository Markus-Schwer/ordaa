use std::{path::Path, future::Future};

use tantivy::{collector::TopDocs, query::QueryParser, schema::*, Index, IndexReader};
use tokio::{
    sync::mpsc::{unbounded_channel, UnboundedReceiver, UnboundedSender},
    task::JoinHandle,
};
use uuid::Uuid;

use super::models::MenuWithItems;

pub struct SearchContextWriter {
    index_write_receiver: UnboundedReceiver<MenuWithItems>,
    index: Index,
}

#[derive(Clone)]
pub struct SearchContextReader {
    pub index_write_sender: UnboundedSender<MenuWithItems>,
    menu_field: Field,
    uuid_field: Field,
    short_name_field: Field,
    name_field: Field,
    parser: QueryParser,
    reader: IndexReader,
}

pub fn init_search_index() -> (SearchContextWriter, SearchContextReader) {
    let (index_write_sender, index_write_receiver) = unbounded_channel::<MenuWithItems>();

    let mut schema_builder = Schema::builder();
    let id_field = schema_builder.add_text_field("uuid", STORED);
    let menu_field = schema_builder.add_text_field("menu", TEXT | STORED);
    let name_field = schema_builder.add_text_field("name", TEXT | STORED);
    let short_name_field = schema_builder.add_text_field("short_name", TEXT | STORED);
    let index = match Index::open_in_dir(Path::new("./index")) {
        Ok(index) => index,
        Err(_) => Index::create_in_dir(Path::new("./index"), schema_builder.build()).unwrap(),
    };
    let writer = SearchContextWriter {
        index: index.clone(),
        index_write_receiver,
    };
    let reader = SearchContextReader {
        menu_field,
        uuid_field: id_field,
        name_field,
        short_name_field,
        index_write_sender,
        parser: QueryParser::for_index(&index, vec![id_field, name_field]),
        reader: index.reader().unwrap(),
    };
    return (writer, reader);
}

impl SearchContextWriter {
    pub fn start_index_writer(mut self, reader: SearchContextReader) -> impl Future<Output=()> {
        let index = self.index.clone();
        async move {
            let mut writer = index.writer(500_000_000).unwrap();
            while let Some(menu_val) = self.index_write_receiver.recv().await {
                for it in menu_val.items {
                    writer
                        .add_document(tantivy::doc!(
                            reader.menu_field => menu_val.name.clone(),
                            reader.uuid_field => it.uuid.to_string(),
                            reader.name_field => it.name,
                            reader.short_name_field => it.short_name
                        ))
                        .unwrap();
                    writer.commit().unwrap();
                }
            }
        }
    }
}

impl SearchContextReader {
    pub fn fuzz_menu_item_ids(&self, prompt: &str) -> Vec<Uuid> {
        let query = self.parser.parse_query_lenient(prompt);
        let searcher = self.reader.searcher();
        let results = searcher.search(&query.0, &TopDocs::with_limit(10)).unwrap();
        let mut matches = Vec::<Uuid>::new();
        for (_, doc_address) in results {
            let retrieved_doc = searcher.doc(doc_address).unwrap();
            let id = retrieved_doc
                .get_first(self.uuid_field)
                .unwrap()
                .as_text()
                .unwrap();
            matches.push(Uuid::parse_str(id).unwrap());
        }
        return matches;
    }
}
