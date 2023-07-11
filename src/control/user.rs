#[derive(PartialEq, Eq, Hash, Clone)]
pub struct User {
    name: String,
}

impl ToString for User {
    fn to_string(&self) -> String {
        self.name.clone()
    }
}
